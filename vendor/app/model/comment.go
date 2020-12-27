package model

import (
	"app/shared/database"
	"time"
)

// *****************************************************************************
// Comment
// *****************************************************************************

// Comment table contains the information for each comment
type Comment struct {
	ID        uint32    `db:"ID"`
	Message   string    `db:"Message"`
	UserId    uint32    `db:"UserId"`
	CreatedAt time.Time `db:"CreatedAt"`
}

// CreateComment creates new comment in the db
func CreateComment(Message string, UserId uint32) (uint32, error) {
	var err error
	var commentID uint32

	err = database.SQL.QueryRow(
		"INSERT INTO comment_table (message, userid, createdat) VALUES ($1, $2, $3) RETURNING id",
		Message, UserId, time.Now()).Scan(&commentID)

	return commentID, StandardizeError(err)
}

// SelectComment gets comment by ID
func SelectComment(id uint32) (Comment, error) {
	var err error

	result := Comment{}

	err = database.SQL.QueryRow("SELECT id, message, userid, createdat FROM comment_table WHERE id = $1 LIMIT 1", id).Scan(
		&result.ID,
		&result.Message,
		&result.UserId,
		&result.CreatedAt)

	return result, StandardizeError(err)
}

// DeleteComment delete comment by ID
func DeleteComment(id uint32) error {
	var err error

	_, err = database.SQL.Exec("DELETE FROM comment_table WHERE id = $1",
		id)

	return StandardizeError(err)
}
