package model

import (
	"app/shared/database"
	"database/sql"
	"encoding/json"
	"time"
)

type ReadBy struct {
	UserId1 uint32 `json:"userId1"`
	UserId2 uint32 `json:"userId2"`
}

// *****************************************************************************
// Chat
// *****************************************************************************

// Chat table contains the information for each chat
type Chat struct {
	ID        uint32          `db:"ID"`
	CreatedAt time.Time       `db:"CreatedAt"`
	Message   string          `db:"Message"`
	ChatId    uint32          `db:"ChatId"`
	ReadBy    json.RawMessage `db:"ReadBy"`
	SendBy    uint32          `db:"SendBy"`
}

// CreateNewMessage creates new message
func CreateNewMessage(message string, userId, chatId uint32) (int, error) {
	var err error
	var messageID int

	readedBy := ReadBy{0, 0}

	readByJson, err := json.Marshal(readedBy)
	if err != nil {
		return 0, StandardizeError(err)
	}

	err = database.SQL.QueryRow(
		"INSERT INTO chat_table (CreatedAt, Message, ChatId, ReadBy, SendBy) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		time.Now(), message, chatId, readByJson, userId).Scan(&messageID)
	if err != nil {
		return 0, StandardizeError(err)
	}

	return messageID, StandardizeError(err)
}

// SelectChat gets all the info of this chat
func SelectChat(chatId, userId uint32) ([]Chat, error) {
	var err error
	var rows *sql.Rows
	var results []Chat
	var result Chat
	var readedBy ReadBy

	rows, err = database.SQL.Query("SELECT * FROM chat_table WHERE chatid = $1 ORDER BY createdat", chatId)
	if err != nil {
		return results, StandardizeError(err)
	}

	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.CreatedAt,
			&result.Message,
			&result.ChatId,
			&result.ReadBy,
			&result.SendBy)
		if err != nil {
			return results, StandardizeError(err)
		}

		// If the message is not read by the user we will set it
		err = json.Unmarshal(result.ReadBy, &readedBy)
		if err != nil {
			return results, StandardizeError(err)
		}
		if result.SendBy != userId && (readedBy.UserId1 == 0 || readedBy.UserId2 == 0) {
			if readedBy.UserId1 == 0 {
				readedBy.UserId1 = userId
			} else if readedBy.UserId1 != userId && readedBy.UserId2 == 0 {
				readedBy.UserId2 = userId
			}
			readByJson, err := json.Marshal(readedBy)
			if err != nil {
				return results, StandardizeError(err)
			}
			_, err = database.SQL.Exec("UPDATE chat_table SET readby = $2 WHERE id = $1", result.ID, readByJson)
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		//TODO:Andrey need fix
		return results, StandardizeError(err)
	}
	return results, StandardizeError(err)
}

// SelectProjectById gets project information by ID
func WhoRead(id uint32) (json.RawMessage, error) {
	var err error
	var result json.RawMessage

	err = database.SQL.QueryRow("SELECT readby FROM chat_table WHERE id=$1", id).Scan(
		&result)

	return result, StandardizeError(err)
}
