// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"database/sql"
	"time"
)

type Contact struct {
	ContactID      int64  `json:"contact_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePhoto   []byte `json:"profile_photo"`
	PhoneNumber    string `json:"phone_number"`
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
}

type GroupMember struct {
	ContactID      int64        `json:"contact_id"`
	GroupID        int64        `json:"group_id"`
	JoinedDatetime time.Time    `json:"joined_datetime"`
	LeftDatetime   sql.NullTime `json:"left_datetime"`
}

type Message struct {
	MessageID    int64     `json:"message_id"`
	Username     string    `json:"username"`
	MessageText  string    `json:"message_text"`
	SentDatetime time.Time `json:"sent_datetime"`
	GroupID      int64     `json:"group_id"`
}

type MessageGroup struct {
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
}
