package model

import (
	"app/shared/database"
	"database/sql"
	"encoding/json"
	"time"
)

// *****************************************************************************
// ProjectArchive
// *****************************************************************************

// ProjectArchive table contains the information for each project archive
type ProjectArchive struct {
	ID                uint32          `db:"ID"`
	CreatedAt         time.Time       `db:"CreatedAt"`
	ApprovedToPresent bool            `db:"ApprovedToPresent"`
	ApprovedFiles     json.RawMessage `db:"ApprovedFiles"`
	ProjectId         uint32          `db:"ProjectId"`
}

// CreateArchiveProject add project with status done to the archive
func CreateArchiveProject(project Project) error {
	var err error

	_, err = database.SQL.Exec(
		"INSERT INTO projectarchive_table (createdat, approvedtopresent, approvedfiles, projectid) VALUES ($1, $2, $3, $4)",
		time.Now(), true, project.Files, project.ID)

	return StandardizeError(err)
}

// ListArchiveProject returns the whole list of all archive projects
func ListArchiveProject() ([]ProjectArchive, error) {
	var err error
	var rows *sql.Rows
	var results []ProjectArchive
	var result ProjectArchive

	rows, err = database.SQL.Query("SELECT id, createdat, approvedtopresent, approvedfiles, projectid FROM projectarchive_table ")
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.CreatedAt,
			&result.ApprovedToPresent,
			&result.ApprovedFiles,
			&result.ProjectId)
		if err != nil {
			return results, StandardizeError(err)
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return results, StandardizeError(err)
	}
	return results, StandardizeError(err)
}

// SelectArchiveProjectById gets archive project information by ID
func SelectArchiveProjectById(id uint32) (ProjectArchive, error) {
	var err error

	result := ProjectArchive{}

	err = database.SQL.QueryRow("SELECT id, createdat, approvedtopresent, approvedfiles, projectid FROM projectarchive_table WHERE id = $1 LIMIT 1", id).Scan(
		&result.ID,
		&result.CreatedAt,
		&result.ApprovedToPresent,
		&result.ApprovedFiles,
		&result.ProjectId,
	)

	return result, StandardizeError(err)
}

// SearchArchiveProject search for archive project
func SearchArchiveProject(query string) ([]ProjectArchive, error) {
	var err error
	var rows *sql.Rows
	var results []ProjectArchive
	var result ProjectArchive

	rows, err = database.SQL.Query("SELECT a.id, a.createdat, a.approvedtopresent, a.approvedfiles, a.projectid FROM projectarchive_table AS a "+
		"JOIN project_table AS t on a.projectid = t.ID "+
		"JOIN user_table AS u1 on (t.StudentsId ->> 'studentId1')::INT = u1.ID "+
		"JOIN user_table AS u2 on (t.StudentsId ->> 'studentId2')::INT = u2.ID "+
		"JOIN user_table AS u3 on  t.supervisorid = u3.ID "+
		"where t.statusid=4 and (t.projectname like '%' || $1 || '%' or u1.firstname like '%'  || $1 || '%' or u1.lastname like '%' || $1 ||'%' or u2.firstname like '%'  || $1 || '%' or u2.lastname like '%' || $1 ||'%' or u3.firstname like '%'  || $1 || '%' or u3.lastname like '%' || $1 ||'%')", query)
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.CreatedAt,
			&result.ApprovedToPresent,
			&result.ApprovedFiles,
			&result.ProjectId,
		)
		if err != nil {
			return results, StandardizeError(err)
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return results, StandardizeError(err)
	}
	return results, StandardizeError(err)
}
