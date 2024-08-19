package api

import (
	db "BreezeServer/db/sqlc"
	"BreezeServer/util"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

type createAccountRequest struct {
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	ProfilePhoto []byte `json:"profile_photo"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required,min=6"`
}

type contactResponse struct {
	FirstName    string `json:"first_name" `
	LastName     string `json:"last_name" `
	ProfilePhoto []byte `json:"profile_photo"`
	PhoneNumber  string `json:"phone_number" `
	Username     string `json:"username" `
	Contact_id   int64  `json:"contact_id"`
}

func newContactResponse(contact db.Contact) contactResponse {
	return contactResponse{
		Username:     contact.Username,
		FirstName:    contact.FirstName,
		LastName:     contact.LastName,
		ProfilePhoto: contact.ProfilePhoto,
		PhoneNumber:  contact.PhoneNumber,
		Contact_id:   contact.ContactID,
	}
}

func (server *Server) createContact(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateContactParams{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		ProfilePhoto:   req.ProfilePhoto,
		PhoneNumber:    req.PhoneNumber,
		Username:       req.Username,
		HashedPassword: hashedPassword,
	}

	contact, err := server.store.CreateContact(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newContactResponse(contact)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string          `json:"access_token"`
	User        contactResponse `json:"user"`
}

func (server *Server) loginContact(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetContact(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newContactResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}

type getContactListRequest struct {
	ContactId int64 `json:"contact_id" binding:"required"`
}

func (server *Server) getContactList(ctx *gin.Context) {
	var req getContactListRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contactList, err := server.store.GetContactList(ctx, req.ContactId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
	}

	ctx.JSON(http.StatusOK, contactList)
}

type getContactListSearchRequest struct {
	Username string `json:"username" binding:"required"`
}

func (server *Server) getContactSearchList(ctx *gin.Context) {
	var req getContactListSearchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contactList, err := server.store.SearchContact(ctx, req.Username)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
	}

	ctx.JSON(http.StatusOK, contactList)
}
