package api

import (
	db "BreezeServer/db/sqlc"
	"BreezeServer/token"
	"context"
	aes2 "crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"io"
	"net/http"
	"sort"
)

type createMessageRequest struct {
	Username    string `json:"username" binding:"required"`
	MessageText string `json:"message_text" binding:"required"`
	GroupId     int64  `json:"group_id" binding:"required"`
}

type listUserGroupMessageRequest struct {
	GroupId int64 `json:"group_id" binding:"required"`
}

type getLastMessage struct {
	GroupID int64 `json:"group_id" binding:"required"`
}

type DeleteMessages struct {
	GroupID int64 `json:"group_id" binding:"required"`
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

func (server *Server) encryptMessage(message string) (string, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(server.config.EncryptionKey)
	if err != nil {
		panic(err)
	}

	if len(decodedKey) != 32 {
		panic(err)
	}

	aes, err := aes2.NewCipher(decodedKey)

	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(message), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (server *Server) decryptMessage(encryptedMessage string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedMessage)

	if err != nil {
		panic(err)
	}

	key, err := base64.StdEncoding.DecodeString(server.config.EncryptionKey)
	if err != nil {
		panic(err)
	}

	block, err := aes2.NewCipher(key)

	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		panic(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}
	return string(plaintext), nil
}

func (server *Server) createMessageLogic(username string, messageText string, groupid int64) (db.Message, error) {
	arg := db.CreateMessageParams{
		Username:    username,
		MessageText: messageText,
		GroupID:     groupid,
	}

	arg.MessageText, _ = server.encryptMessage(arg.MessageText)

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

	messages, err := server.store.ListUserGroupMessage(ctx, req.GroupId)
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

	sort.Slice(messages[:], func(i, j int) bool {
		return messages[i].SentDatetime.Before(messages[j].SentDatetime)
	})
	for i := range messages {
		messages[i].MessageText, _ = server.decryptMessage(messages[i].MessageText)
	}

	ctx.JSON(http.StatusOK, messages)
}

func (server *Server) getLastMessage(ctx *gin.Context) {
	var req getLastMessage

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	message, err := server.store.FetchLastMessage(ctx, req.GroupID)
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
	ctx.JSON(http.StatusOK, message)
}

func (server *Server) deleteMessages(ctx *gin.Context) {
	var req DeleteMessages

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteMessages(ctx, req.GroupID)

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

	ctx.JSON(http.StatusOK, "Messages delete successfully.")
}
