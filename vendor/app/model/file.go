package model

import (
	"app/shared/database"
)

// *****************************************************************************
// File
// *****************************************************************************

// File table contains the information for each file
type File struct {
	ID   uint32 `db:"ID"`
	Link string `db:"Link"`
	Type uint32 `db:"Type"`
}

// CreateFile creates file
func CreateFile(link string, Type uint32) (uint32, error) {
	var err error
	var id uint32

	err = database.SQL.QueryRow("INSERT INTO file_table (link, type) VALUES ($1, $2) RETURNING id",
		link, Type).Scan(&id)

	return id, StandardizeError(err)
}

// SelectFile gets file information by ID
func SelectFile(ID uint32) (File, error) {
	var err error

	result := File{}

	err = database.SQL.QueryRow("SELECT * FROM file_table WHERE id = $1 LIMIT 1", ID).Scan(
		&result.ID,
		&result.Link,
		&result.Type)

	return result, StandardizeError(err)
}
