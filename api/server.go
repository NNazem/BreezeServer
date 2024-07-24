package api

import (
	db "BreezeServer/db/sqlc"
	"BreezeServer/token"
	"BreezeServer/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
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
	}

	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	fmt.Println("Setting up router")
	router := gin.Default()

	router.POST("/contacts", server.createContact)
	router.POST("/contacts/login", server.loginContact)

	//authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	//auth

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
