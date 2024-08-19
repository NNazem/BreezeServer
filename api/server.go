package api

import (
	db "BreezeServer/db/sqlc"
	"BreezeServer/token"
	"BreezeServer/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Server struct {
	config     util.Config
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	room       *room
	rooms      map[int64]*room
	roomsMutex sync.RWMutex
}

func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		room:       newRoom(),
		rooms:      make(map[int64]*room),
	}

	server.room.setServer(server)

	server.setUpRouter()
	return server, nil
}

func (server *Server) getOrCreateRoom(roomID int64) *room {
	server.roomsMutex.Lock()
	defer server.roomsMutex.Unlock()

	if r, exists := server.rooms[roomID]; exists {
		return r
	}

	r := newRoom()
	r.setServer(server)
	go r.run()
	server.rooms[roomID] = r
	log.Printf("Nuova stanza creata con ID: %s", roomID)
	return r
}

func (server *Server) setUpRouter() {
	fmt.Println("Setting up router")
	router := gin.Default()

	router.Use(corsMiddleware())
	router.POST("/contacts", server.createContact)
	router.POST("/contacts/login", server.loginContact)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/contacts/getContactLists", server.getContactList)
	authRoutes.POST("/contacts/search", server.getContactSearchList)
	authRoutes.POST("/messageGroups", server.createMessageGroup)
	authRoutes.POST("/groupMembers", server.createGroupMember)
	authRoutes.POST("/groupMembers/searchId", server.getGroupId)

	authRoutes.POST("/messages", server.createMessage)
	authRoutes.POST("/messages/fetchMessages", server.listUserGroupMessage)
	authRoutes.POST("/messages/last", server.getLastMessage)

	router.GET("/ws", func(c *gin.Context) {
		roomID := c.Query("group_id")
		token := c.Query("token")
		username := c.Query("username")
		payload, err := server.tokenMaker.VerifyToken(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if payload.Username != username {
			err := errors.New("Username doesn't belong to the authenticated user")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		roomIDConverted, _ := strconv.ParseInt(roomID, 10, 64)
		room := server.getOrCreateRoom(roomIDConverted)
		server.serveWs(c.Writer, c.Request, room)
	})

	server.router = router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) serveWs(w http.ResponseWriter, r *http.Request, room *room) {
	room.ServerHTTP(w, r)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
