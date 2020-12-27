package model

import (
	"app/shared/database"
)

// *****************************************************************************
// ApprovalForm
// *****************************************************************************

// ApprovalForm table contains the information for each approval form
type ApprovalForm struct {
	ID                uint32 `db:"ID"`
	Synopses          string `db:"Synopses"`
	ScopeOfTheProject string `db:"ScopeOfTheProject"`
	UniqueFeatures    string `db:"UniqueFeatures"`
}

// CreateApprovalForm creates new approval form in the db and returns ID
func CreateApprovalForm(synopses, scopeOfTheProject, uniqueFeatures string, UserId uint32) (uint32, error) {
	var err error
	var approvalFormID uint32

	err = database.SQL.QueryRow(
		"INSERT INTO approvalform_table (synopses, scopeoftheproject, uniquefeatures) VALUES ($1, $2, $3) RETURNING id",
		synopses, scopeOfTheProject, uniqueFeatures).Scan(&approvalFormID)

	return approvalFormID, StandardizeError(err)
}

// SelectApprovalForm gets approval form by ID
func SelectApprovalForm(id uint32) (ApprovalForm, error) {
	var err error

	result := ApprovalForm{}

	err = database.SQL.QueryRow("SELECT id, synopses, scopeoftheproject, uniquefeatures FROM approvalform_table WHERE id = $1 LIMIT 1", id).Scan(
		&result.ID,
		&result.Synopses,
		&result.ScopeOfTheProject,
		&result.UniqueFeatures)

	return result, StandardizeError(err)
}

// DeleteApprovalForm delete approval form by ID
func DeleteApprovalForm(id uint32) error {
	var err error

	_, err = database.SQL.Exec("DELETE FROM approvalform_table WHERE id = $1",
		id)

	return StandardizeError(err)
}
