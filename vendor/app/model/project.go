package model

import (
	"app/shared/database"
	"database/sql"
	"encoding/json"
	"time"
)

// *****************************************************************************
// Project
// *****************************************************************************

// Project table contains the information for each project
type Project struct {
	ID               uint32          `db:"ID"`
	ProjectName      string          `db:"ProjectName"`
	Description      string          `db:"Description"`
	ShortDescription string          `db:"ShortDescription"`
	StatusId         uint32          `db:"StatusId"`
	Files            json.RawMessage `db:"Files"`
	Type             uint            `db:"Type"`
	FormId           uint32          `db:"FormId"`
	CreateAt         time.Time       `db:"CreateAt"`
	UpdateAt         time.Time       `db:"UpdateAt"`
	CommentsId       json.RawMessage `db:"CommentsId"`
	StudentsId       json.RawMessage `db:"StudentsId"`
	SupervisorId     uint32          `db:"SupervisorId"`
}

type Students struct {
	StudentId1 int `json:"studentId1"`
	StudentId2 int `json:"studentId2"`
}

type ProjectFiles struct {
	Book1PDF      int `json:"Book1PDF"`
	Book1WORD     int `json:"Book1WORD"`
	Presentation1 int `json:"Presentation1"`
	Book2PDF      int `json:"Book2PDF"`
	Book2WORD     int `json:"Book2WORD"`
	Presentation2 int `json:"Presentation2"`
	SourceCode    int `json:"SourceCode"`
}

// SelectProjectById gets project information by ID
func SelectProjectById(id uint32) (Project, error) {
	var err error

	result := Project{}

	err = database.SQL.QueryRow("SELECT id, projectname, description, shortdescription, statusid, type, formid, createdat, updateat, supervisorid, studentsid,commentsid FROM project_table WHERE id = $1 LIMIT 1", id).Scan(
		&result.ID,
		&result.ProjectName,
		&result.Description,
		&result.ShortDescription,
		&result.StatusId,
		&result.Type,
		&result.FormId,
		&result.CreateAt,
		&result.UpdateAt,
		&result.SupervisorId,
		&result.StudentsId,
		&result.CommentsId)

	return result, StandardizeError(err)
}

// SelectProjectByStudent checks if the student assign to project. If yes returns the project
func SelectProjectByStudent(studentId uint32) (Project, error) {
	var err error

	result := Project{}

	err = database.SQL.QueryRow("SELECT id, projectname, description, shortdescription, statusid, type, formid, createdat, updateat, supervisorid, studentsid, commentsid, files FROM project_table WHERE (StudentsId ->> 'studentId1')::INTEGER = $1 or (StudentsId ->> 'studentId2')::INTEGER = $1 LIMIT 1", studentId).Scan(
		&result.ID,
		&result.ProjectName,
		&result.Description,
		&result.ShortDescription,
		&result.StatusId,
		&result.Type,
		&result.FormId,
		&result.CreateAt,
		&result.UpdateAt,
		&result.SupervisorId,
		&result.StudentsId,
		&result.CommentsId,
		&result.Files)

	return result, StandardizeError(err)
}

// CreateStudentIdea creates student idea
func CreateStudentIdea(ProjectName, Description, ShortDescription string, SupervisorId uint32, StudentsId json.RawMessage) error {
	var err error

	_, err = database.SQL.Exec(
		"INSERT INTO project_table (ProjectName, Description, ShortDescription, StatusId, Type, formid, CreatedAt, UpdateAt, StudentsId, SupervisorId, commentsid, files) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		ProjectName, Description, ShortDescription, 2, 1, 0, time.Now(), time.Now(), StudentsId, SupervisorId, "{}", "{}")

	return StandardizeError(err)
}

// CreateProject creates project by the supervisor
func CreateProject(ProjectName, Description, ShortDescription string, SupervisorId uint32, StudentsId json.RawMessage) error {
	var err error

	_, err = database.SQL.Exec(
		"INSERT INTO project_table (ProjectName, Description, ShortDescription, StatusId, Type, formid, CreatedAt, UpdateAt, StudentsId, SupervisorId, commentsid, files) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		ProjectName, Description, ShortDescription, 1, 2, 0, time.Now(), time.Now(), StudentsId, SupervisorId, "{}", "{}")

	return StandardizeError(err)
}

// ListOfStudentsWithoutProject return the list of user who are not assigned to any project
func ListOfStudentsWithoutProject(id uint32) ([]User, error) {
	var err error
	var rows *sql.Rows
	var results []User
	var result User

	rows, err = database.SQL.Query("SELECT * FROM user_table "+
		"WHERE "+
		"user_table.firstname is not NULL and "+
		"user_table.lastname is not NULL and "+
		"block is FALSE and "+
		"id != $1 and "+
		"id not in (SELECT d.value::INTEGER FROM project_table JOIN json_each_text(project_table.StudentsId) d ON true WHERE d.value::INTEGER != 0) and "+
		"role = 1 order by 1", id)
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

// ListOpenProjects returns a list of all projects with status ID open
func ListOpenProjects() ([]Project, error) {
	var err error
	var rows *sql.Rows
	var results []Project
	var result Project

	rows, err = database.SQL.Query("SELECT id, projectname, shortdescription, supervisorid FROM project_table WHERE statusid=1")
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.ProjectName,
			&result.ShortDescription,
			&result.SupervisorId)
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

// ListWaitProjects returns the list of the project which waiting for the supervisor approval
func ListWaitProjects(id uint32) ([]Project, error) {
	var err error
	var rows *sql.Rows
	var results []Project
	var result Project

	rows, err = database.SQL.Query("SELECT id, projectname, shortdescription, studentsid, type FROM project_table WHERE statusid=2 and supervisorid=$1", id)
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.ProjectName,
			&result.ShortDescription,
			&result.StudentsId,
			&result.Type)
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

// ListWaitManagerApproval returns the list of the project which waiting for the project manager approval
func ListWaitManagerApproval() ([]Project, error) {
	var err error
	var rows *sql.Rows
	var results []Project
	var result Project

	rows, err = database.SQL.Query("SELECT id, projectname, shortdescription, studentsid, type, formid FROM project_table WHERE statusid=3")
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.ProjectName,
			&result.ShortDescription,
			&result.StudentsId,
			&result.Type,
			&result.FormId)
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

// ListRunningProjects returns a list of all projects have been approved
func ListRunningProjects() ([]Project, error) {
	var err error
	var rows *sql.Rows
	var results []Project
	var result Project

	rows, err = database.SQL.Query("SELECT id, projectname, shortdescription, studentsid, type FROM project_table WHERE statusid=4")
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.ProjectName,
			&result.ShortDescription,
			&result.StudentsId,
			&result.Type)
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

// ListRunningProjects returns a list of all projects have been approved
func ListRunningBySupervisorProjects(id uint32) ([]Project, error) {
	var err error
	var rows *sql.Rows
	var results []Project
	var result Project

	rows, err = database.SQL.Query("SELECT id, projectname, shortdescription, studentsid, type FROM project_table WHERE statusid=4 and supervisorid=$1", id)
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.ProjectName,
			&result.ShortDescription,
			&result.StudentsId,
			&result.Type)
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

// AssignStudentsToProject assigned two users to the project
func AssignStudentsToProject(id uint32, StudentsId json.RawMessage) error {
	var err error

	_, err = database.SQL.Exec("UPDATE project_table SET studentsid=$2, statusid=2, updateat=$3 WHERE id = $1",
		id, StudentsId, time.Now())

	return StandardizeError(err)
}

// UpdateComments the array of comments with the new comment
func UpdateComments(id uint32, comments json.RawMessage) error {
	var err error

	_, err = database.SQL.Exec("UPDATE project_table SET commentsid=$2, updateat=$3 WHERE id = $1",
		id, comments, time.Now())

	return StandardizeError(err)
}

// UpdateProject update the name, description and short description of the project
func UpdateProject(id uint32, name, description, shortDescription string) error {
	var err error

	_, err = database.SQL.Exec("UPDATE project_table SET projectname=$2, description=$3, shortdescription=$4, updateat=$5 WHERE id = $1",
		id, name, description, shortDescription, time.Now())

	return StandardizeError(err)
}

// ApproveProject supervisor approves the project by ID
func ApproveProject(id uint32) error {
	var err error

	_, err = database.SQL.Exec("UPDATE project_table SET statusid=3, updateat=$2 WHERE id = $1",
		id, time.Now())

	return StandardizeError(err)
}

// ApproveProject supervisor approves the project by ID
func ApproveProjectByManager(id uint32) error {
	var err error

	_, err = database.SQL.Exec("UPDATE project_table SET statusid=4, updateat=$2 WHERE id = $1",
		id, time.Now())

	return StandardizeError(err)
}

// DeclineProject supervisor decline the project by ID
func DeclineProject(id uint32) error {
	var err error

	_, err = database.SQL.Exec("UPDATE project_table SET statusid=1, studentsid='{\"studentId1\":0,\"studentId2\":0}', updateat=$2 WHERE id = $1",
		id, time.Now())

	return StandardizeError(err)
}

// DeleteProject delete project or idea from the data base by the ID
func DeleteProject(id uint32) error {
	var err error

	_, err = database.SQL.Exec("DELETE FROM project_table WHERE id = $1",
		id)

	return StandardizeError(err)
}

// UpdateComments the array of comments with the new comment
func SetApprovalForm(id, formId uint32) error {
	var err error

	_, err = database.SQL.Exec("UPDATE project_table SET formid=$2, updateat=$3 WHERE id = $1",
		id, formId, time.Now())

	return StandardizeError(err)
}

// SearchOpenProject search for project with status open
func SearchOpenProject(query string) ([]Project, error) {
	var err error
	var rows *sql.Rows
	var results []Project
	var result Project

	rows, err = database.SQL.Query("SELECT t.id, t.projectname, t.shortdescription, t.supervisorid FROM project_table as t left join user_table as u on t.supervisorid = u.id WHERE t.statusid=1 and (t.projectname like '%' || $1 || '%' or u.firstname like '%'  || $1 || '%' or u.lastname like '%' || $1 ||'%')", query)
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.ProjectName,
			&result.ShortDescription,
			&result.SupervisorId)
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

// SearchWaitingProject search for project with status waiting
func SearchWaitingProject(id uint32, query string) ([]Project, error) {
	var err error
	var rows *sql.Rows
	var results []Project
	var result Project

	rows, err = database.SQL.Query("SELECT t.id, t.projectname, t.shortdescription, t.studentsid, t.type FROM project_table AS t "+
		"JOIN user_table AS u1 on (t.StudentsId ->> 'studentId1')::INT = u1.ID "+
		"JOIN user_table AS u2 on (t.StudentsId ->> 'studentId2')::INT = u2.ID "+
		"where t.statusid=2 and t.supervisorid=$2 and (t.projectname like '%' || $1 || '%' or u1.firstname like '%'  || $1 || '%' or u1.lastname like '%' || $1 ||'%' or u2.firstname like '%'  || $1 || '%' or u2.lastname like '%' || $1 ||'%')", query, id)
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.ProjectName,
			&result.ShortDescription,
			&result.StudentsId,
			&result.Type)
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

// SearchManagerApprovalProject search for project with status waiting for project manager approval
func SearchManagerApprovalProject(query string) ([]Project, error) {
	var err error
	var rows *sql.Rows
	var results []Project
	var result Project

	rows, err = database.SQL.Query("SELECT t.id, t.projectname, t.shortdescription, t.studentsid, t.type FROM project_table AS t "+
		"JOIN user_table AS u1 on (t.StudentsId ->> 'studentId1')::INT = u1.ID "+
		"JOIN user_table AS u2 on (t.StudentsId ->> 'studentId2')::INT = u2.ID "+
		"JOIN user_table AS u3 on  t.supervisorid = u3.ID "+
		"where t.statusid=3 and (t.projectname like '%' || $1 || '%' or u1.firstname like '%'  || $1 || '%' or u1.lastname like '%' || $1 ||'%' or u2.firstname like '%'  || $1 || '%' or u2.lastname like '%' || $1 ||'%' or u3.firstname like '%'  || $1 || '%' or u3.lastname like '%' || $1 ||'%')", query)
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.ProjectName,
			&result.ShortDescription,
			&result.StudentsId,
			&result.Type)
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

// SearchRunningProjectBySupervisor search for project with status running by supervisor id
func SearchRunningProjectBySupervisor(id uint32, query string) ([]Project, error) {
	var err error
	var rows *sql.Rows
	var results []Project
	var result Project

	rows, err = database.SQL.Query("SELECT t.id, t.projectname, t.shortdescription, t.studentsid, t.type FROM project_table AS t "+
		"JOIN user_table AS u1 on (t.StudentsId ->> 'studentId1')::INT = u1.ID "+
		"JOIN user_table AS u2 on (t.StudentsId ->> 'studentId2')::INT = u2.ID "+
		"where t.statusid=4 and t.supervisorid=$2 and (t.projectname like '%' || $1 || '%' or u1.firstname like '%'  || $1 || '%' or u1.lastname like '%' || $1 ||'%' or u2.firstname like '%'  || $1 || '%' or u2.lastname like '%' || $1 ||'%')", query, id)
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.ProjectName,
			&result.ShortDescription,
			&result.StudentsId,
			&result.Type)
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

// SearchRunningProject search for project with status running
func SearchRunningProject(query string) ([]Project, error) {
	var err error
	var rows *sql.Rows
	var results []Project
	var result Project

	rows, err = database.SQL.Query("SELECT t.id, t.projectname, t.shortdescription, t.studentsid, t.type FROM project_table AS t "+
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
			&result.ProjectName,
			&result.ShortDescription,
			&result.StudentsId,
			&result.Type)
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

// AddFileToListOfFiles add new file to the list of files
func AddFileToListOfFiles(id uint32, files json.RawMessage) error {
	var err error

	_, err = database.SQL.Exec("UPDATE project_table SET files=$2, updateat=$3 WHERE id = $1",
		id, files, time.Now())

	return StandardizeError(err)
}

// MarkProjectAsDone the project is DONE!!!
func MarkProjectAsDone(id uint32) error {
	var err error

	_, err = database.SQL.Exec("UPDATE project_table SET statusid=5, updateat=$2 WHERE id = $1",
		id, time.Now())

	return StandardizeError(err)
}
