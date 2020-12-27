package model

import (
	"app/shared/database"
	"database/sql"
)

// *****************************************************************************
// ProgressBar
// *****************************************************************************

// ProgressBar table contains the information for each progress bar
type ProgressBar struct {
	ID            uint32 `db:"ID"`
	ProgressBarId uint32 `db:"ProgressBarId"`
	Milestone     string `db:"Milestone"`
	Done          bool   `db:"Done"`
}

// CreateMilestone add new milestone to the database
func CreateMilestone(progressBarId uint32, milestone string) error {
	var err error

	_, err = database.SQL.Exec(
		"INSERT INTO progressbar_table (progressbarid, milestone, done) VALUES ($1, $2, $3)",
		progressBarId, milestone, false)

	return StandardizeError(err)
}

// MarkAsDone mark the milestone as done
func MarkAsDone(ID uint32) error {
	var err error

	_, err = database.SQL.Exec("UPDATE progressbar_table SET Done=$2 WHERE id = $1", ID, true)

	return StandardizeError(err)
}

// DeleteMilestone delete the milestone
func DeleteMilestone(ID uint32) error {
	var err error

	_, err = database.SQL.Exec("DELETE FROM progressbar_table WHERE id = $1", ID)

	return StandardizeError(err)
}

// SelectListOfMilestones gets all milestones
func SelectListOfMilestones(progressBarId uint32) ([]ProgressBar, error) {
	var err error
	var rows *sql.Rows
	var results []ProgressBar
	var result ProgressBar

	rows, err = database.SQL.Query("SELECT * FROM progressbar_table WHERE ProgressBarId = $1 ORDER BY 1", progressBarId)
	if err != nil {
		return results, StandardizeError(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.ProgressBarId,
			&result.Milestone,
			&result.Done)
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
