package api

import (
	db "BreezeServer/db/sqlc"
	"BreezeServer/token"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

type createMessageRequest struct {
	Username    string `json:"username" binding:"required"`
	MessageText string `json:"message_text" binding:"required"`
	GroupId     int64  `json:"group_id" binding:"required"`
}

type listUserGroupMessageRequest struct {
	Username string `json:"username" binding:"required"`
	GroupId  int64  `json:"group_id" binding:"required"`
}

func (server *Server) createMessage(ctx *gin.Context) {
	var req createMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if req.Username != authPayload.Username {
		err := errors.New("Username doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	message, err := server.createMessageLogic(req.Username, req.MessageText, req.GroupId)

	if err != nil {
		if err.Error() == "forbidden action" {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}
	ctx.JSON(http.StatusOK, message)
}

func (server *Server) createMessageLogic(username string, messageText string, groupid int64) (db.Message, error) {
	arg := db.CreateMessageParams{
		Username:    username,
		MessageText: messageText,
		GroupID:     groupid,
	}

	message, err := server.store.CreateMessage(context.Background(), arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foregin_key_violation", "unique_violation":
				return db.Message{}, errors.New("forbidden action")
			}
		}
		return db.Message{}, err
	}

	return message, nil
}

func (server *Server) listUserGroupMessage(ctx *gin.Context) {
	var req listUserGroupMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUserGroupMessageParams{
		Username: req.Username,
		GroupID:  req.GroupId,
	}

	messages, err := server.store.ListUserGroupMessage(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messages)
}
