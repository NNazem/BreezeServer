package api

import (
	db "BreezeServer/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

type createMessageGroupRequest struct {
	GroupName string `json:"group_name" binding:"required"`
}

type messageGroupResponse struct {
	GroupId   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
}

func newMessageGroupResponse(messageGroup db.MessageGroup) messageGroupResponse {
	return messageGroupResponse{
		GroupId:   messageGroup.GroupID,
		GroupName: messageGroup.GroupName,
	}
}

func (server *Server) createMessageGroup(ctx *gin.Context) {
	var req createMessageGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.CreateMessageGroup(ctx, req.GroupName)

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

	ctx.JSON(http.StatusOK, account)
}
