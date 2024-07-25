package api

import (
	db "BreezeServer/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

type createGroupMemberRequest struct {
	ContactId int64 `json:"contact_id" binding:"required"`
	GroupId   int64 `json:"group_id" binding:"required"`
}

func (server *Server) createGroupMember(ctx *gin.Context) {
	var req createGroupMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateGroupMemberParams{
		ContactID: req.ContactId,
		GroupID:   req.GroupId,
	}

	groupMember, err := server.store.CreateGroupMember(ctx, arg)
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

	ctx.JSON(http.StatusOK, groupMember)
}
