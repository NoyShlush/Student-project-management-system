package model

import (
	"app/shared/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type User struct {
	ID                uint32          `db:"ID"`
	IdNumber          uint32          `db:"IdNumber"`
	FirstName         string          `db:"FirstName"`
	LastName          string          `db:"LastName"`
	Email             string          `db:"Email"`
	Role              uint            `db:"Role"`
	ContactInfomation json.RawMessage `db:"ContactInfomation"`
	HashPassword      string          `db:"HashPassword"`
	CreatedAt         time.Time       `db:"CreatedAt"`
	Block             bool            `db:"Block"`
	ResetPasswordAt   sql.NullTime    `db:"ResetPasswordAt"`
	ResetToken        sql.NullString  `db:"ResetToken"`
	TokenValidUntil   sql.NullTime    `db:"TokenValidUntil"`
}

type ContactInfo struct {
	Telephone int    `json:"Telephone"`
	Facebook  string `json:"Facebook"`
	Twitter   string `json:"Twitter"`
	Telegram  string `json:"Telegram"`
}

// UserID returns the user id
func (u *User) UserID() string {
	r := ""
	r = fmt.Sprintf("%v", u.ID)
	return r
}

// SelectUserByEmail gets user information by email
func SelectUserByEmail(email string) (User, error) {
	var err error

	result := User{}

	err = database.SQL.QueryRow("SELECT * FROM user_table WHERE email = $1 LIMIT 1", email).Scan(
		&result.ID,
		&result.IdNumber,
		&result.FirstName,
		&result.LastName,
		&result.Email,
		&result.Role,
		&result.ContactInfomation,
		&result.HashPassword,
		&result.CreatedAt,
		&result.Block,
		&result.ResetPasswordAt,
		&result.ResetToken,
		&result.TokenValidUntil)

	return result, StandardizeError(err)
}

// SelectUserByIdNumber gets user information by ID number
func SelectUserByIdNumber(idNumber uint32) (User, error) {
	var err error

	result := User{}

	err = database.SQL.QueryRow("SELECT * FROM user_table WHERE idnumber = $1 LIMIT 1", idNumber).Scan(
		&result.ID,
		&result.IdNumber,
		&result.FirstName,
		&result.LastName,
		&result.Email,
		&result.Role,
		&result.ContactInfomation,
		&result.HashPassword,
		&result.CreatedAt,
		&result.Block,
		&result.ResetPasswordAt,
		&result.ResetToken,
		&result.TokenValidUntil)

	return result, StandardizeError(err)
}

// SelectUser gets user information by ID
func SelectProjectManager() (User, error) {
	var err error

	result := User{}

	err = database.SQL.QueryRow("SELECT id, firstname, lastname, idnumber, role, email, block, contactinfomation FROM user_table WHERE role = 3 LIMIT 1").Scan(
		&result.ID,
		&result.FirstName,
		&result.LastName,
		&result.IdNumber,
		&result.Role,
		&result.Email,
		&result.Block,
		&result.ContactInfomation)

	return result, StandardizeError(err)
}

// SelectUser gets user information by ID
func SelectUser(ID uint32) (User, error) {
	var err error

	result := User{}

	err = database.SQL.QueryRow("SELECT id, idnumber, email, block FROM user_table WHERE id = $1 LIMIT 1", ID).Scan(
		&result.ID,
		&result.IdNumber,
		&result.Email,
		&result.Block)

	return result, StandardizeError(err)
}

// CheckIfTheUserExists checking if user existing in the DB by ID number or email
func CheckIfTheUserExists(IdNumber uint32, email string) (User, error) {
	var err error

	result := User{}

	err = database.SQL.QueryRow("SELECT id, idnumber, email FROM user_table WHERE idnumber = $1 or email = $2 LIMIT 1", IdNumber, email).Scan(
		&result.ID,
		&result.IdNumber,
		&result.Email)

	return result, StandardizeError(err)
}

// UserCreate creates user
func CreateNewUser(idNumber uint32, email, password string, role uint) (uint32, error) {
	var err error
	var lastId uint32

	err = database.SQL.QueryRow(
		"INSERT INTO user_table (idnumber, email, role, hashpassword, createdat, block, contactinfomation) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ID",
		idNumber, email, role, password, time.Now(), false, "{}").Scan(&lastId)
	if err != nil {
		return 0, err
	}

	return lastId, StandardizeError(err)
}

// UpdateInfo update the user information
func UpdateInfo(ID uint32, firstName, lastName string, ContactInformation json.RawMessage) error {
	var err error

	_, err = database.SQL.Exec("UPDATE user_table SET firstName=$2, lastName=$3, contactInfomation=$4 WHERE id = $1",
		ID, firstName, lastName, ContactInformation)

	return StandardizeError(err)
}

// SelectUser gets user information by ID
func SelectUserInfo(ID uint32) (User, error) {
	var err error

	result := User{}

	err = database.SQL.QueryRow("SELECT id, firstname, lastname, idnumber, role, email, block, contactinfomation FROM user_table WHERE id = $1 LIMIT 1", ID).Scan(
		&result.ID,
		&result.FirstName,
		&result.LastName,
		&result.IdNumber,
		&result.Role,
		&result.Email,
		&result.Block,
		&result.ContactInfomation)

	return result, StandardizeError(err)
}

// UpdateUser update the user basic data
func UpdateUser(ID, IdNumber uint32, email string) error {
	var err error

	_, err = database.SQL.Exec("UPDATE user_table SET idnumber=$2, email=$3 WHERE id = $1",
		ID, IdNumber, email)

	return StandardizeError(err)
}

// BlockUser block the user
func BlockUser(ID uint32) error {
	var err error

	_, err = database.SQL.Exec("UPDATE user_table SET block = true WHERE id = $1",
		ID)

	return StandardizeError(err)
}

// UnBlockUser unblock the user
func UnBlockUser(ID uint32) error {
	var err error

	_, err = database.SQL.Exec("UPDATE user_table SET block = false WHERE id = $1",
		ID)

	return StandardizeError(err)
}

// SelectListOfUsers gets all user information by role
func SelectListOfUsers(role uint) ([]User, error) {
	var err error
	var rows *sql.Rows
	var results []User
	var result User

	rows, err = database.SQL.Query("SELECT * FROM user_table "+
		"WHERE "+
		"user_table.firstname is not NULL and "+
		"user_table.lastname is not NULL and "+
		"role = $1 order by 1", role)
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.IdNumber,
			&result.FirstName,
			&result.LastName,
			&result.Email,
			&result.Role,
			&result.ContactInfomation,
			&result.HashPassword,
			&result.CreatedAt,
			&result.Block,
			&result.ResetPasswordAt,
			&result.ResetToken,
			&result.TokenValidUntil)
		if err != nil {
			return results, StandardizeError(err)
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		//TODO:Andrey need fix
		return results, StandardizeError(err)
	}
	return results, StandardizeError(err)
}

// CreateChangePasswordToken generates token for change password flow
func CreateChangePasswordToken(ID uint32) (string, error) {
	var err error

	token, err := uuid.NewUUID()
	if err != nil {
		return token.String(), StandardizeError(err)
	}

	_, err = database.SQL.Exec("UPDATE user_table SET resettoken=$2, tokenvaliduntil=$3 WHERE id = $1",
		ID, token, time.Now().Add(time.Hour))

	return token.String(), StandardizeError(err)
}

// UpdateChangePassword update the user password
func UpdateChangePassword(ID uint32, password string) error {
	var err error

	_, err = database.SQL.Exec("UPDATE user_table SET hashpassword=$2 WHERE id = $1",
		ID, password)

	return StandardizeError(err)
}

// CheckUserToken checks if the token exists and returns id if exist or error
func CheckUserToken(token string) (uint32, time.Time, error) {
	var err error
	var id uint32
	var validUntil time.Time

	err = database.SQL.QueryRow("SELECT id, tokenvaliduntil FROM user_table WHERE resetToken = $1 LIMIT 1", token).Scan(
		&id, &validUntil)

	return id, validUntil, StandardizeError(err)
}

// SearchUser search for user in the users table by name, idnumber and email
func SearchUser(query string, role uint) ([]User, error) {
	var err error
	var rows *sql.Rows
	var results []User
	var result User

	rows, err = database.SQL.Query("SELECT id, CAST(idnumber as varchar), firstname, lastname, email, block FROM user_table WHERE role = $2 and firstname is not NULL and lastname is not NULL and (CAST(idnumber as varchar) like '%' || $1 || '%' or  firstname like  '%' || $1 || '%' or lastname like '%' || $1 || '%' or email like '%' || $1 || '%')", query, role)
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.IdNumber,
			&result.FirstName,
			&result.LastName,
			&result.Email,
			&result.Block)
		if err != nil {
			return results, StandardizeError(err)
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		//TODO:Andrey need fix
		return results, StandardizeError(err)
	}
	return results, StandardizeError(err)
}

// SearchWaitingUser search for user in the users table by name, idnumber and email
func SearchWaitingUser(query string, role uint) ([]User, error) {
	var err error
	var rows *sql.Rows
	var results []User
	var result User

	rows, err = database.SQL.Query("SELECT id, CAST(idnumber as varchar), email FROM user_table WHERE role = $2 and (CAST(idnumber as varchar) like '%' || $1 || '%' or  firstname like  '%' || $1 || '%' or lastname like '%' || $1 || '%' or email like '%' || $1 || '%')", query, role)
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.IdNumber,
			&result.Email)
		if err != nil {
			return results, StandardizeError(err)
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		//TODO:Andrey need fix
		return results, StandardizeError(err)
	}
	return results, StandardizeError(err)
}

//SelectListOfWaiting gets all the users that waiting in the users table by role
func SelectListOfWaiting(role uint) ([]User, error) {
	var err error
	var rows *sql.Rows
	var results []User
	var result User

	rows, err = database.SQL.Query("SELECT id, idnumber, email FROM user_table "+
		"WHERE "+
		"user_table.firstname is NULL and "+
		"user_table.lastname is NULL and "+
		"role = $1", role)
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.IdNumber,
			&result.Email)
		if err != nil {
			return results, StandardizeError(err)
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		//TODO:Andrey need fix
		return results, StandardizeError(err)
	}
	return results, StandardizeError(err)
}
