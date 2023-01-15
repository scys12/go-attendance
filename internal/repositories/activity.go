package repositories

import (
	"database/sql"

	"github.com/scys12/go-attendance/internal/models"
	"github.com/scys12/go-attendance/internal/types"
	"github.com/sirupsen/logrus"
)

type ActivityRepo struct {
	db *sql.DB
}

func NewActivityRepo(db *sql.DB) IActivityRepo {
	return &ActivityRepo{
		db: db,
	}
}

func (a *ActivityRepo) AddActivity(activity models.Activity) (int64, error) {
	var id int64
	res, err := a.db.Exec(`
	INSERT INTO activities(
		description,
		attendance_id,
		user_id
	) VALUES (
		?,
		?,
		?
	)`,
		activity.Description,
		activity.AttendanceID,
		activity.UserID)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "AddActivity",
			"error":    err.Error(),
		}).Errorln("[ActivityRepo] Failed to insert activity")
		return 0, types.ErrFailedInsertActivity
	}

	id, err = res.LastInsertId()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "AddActivity",
			"error":    err.Error(),
		}).Errorln("[ActivityRepo] Failed to get ID")

		return 0, types.ErrFailedGetLastIDDB
	}

	return id, nil
}

func (a *ActivityRepo) UpdateActivity(activity models.Activity) error {
	res, err := a.db.Exec(`
	UPDATE activities SET
		description = ?
	WHERE
		id = ? and user_id = ?
	`,
		activity.Description,
		activity.ID,
		activity.UserID,
	)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "UpdateActivity",
			"error":    err.Error(),
		}).Errorln("[ActivityRepo] Failed to update activity")
		return types.ErrFailedUpdateActivity
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "UpdateActivity",
			"error":    err.Error(),
		}).Errorln("[ActivityRepo] Failed to get rows affected")

		return types.ErrFailedRowsAffected
	}

	if rowsAffected < 1 {
		return types.ErrNoRowsAffected
	}

	return nil
}

func (a *ActivityRepo) DeleteActivity(id, userID int64) error {
	res, err := a.db.Exec(`
	DELETE 
	FROM activities
	WHERE
		id = ? and user_id = ?
	`,
		id,
		userID,
	)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "DeleteActivity",
			"error":    err.Error(),
		}).Errorln("[ActivityRepo] Failed to delete activity")

		return types.ErrFailedDeleteActivity
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "DeleteActivity",
			"error":    err.Error(),
		}).Errorln("[ActivityRepo] Failed to get rows affected")

		return types.ErrFailedRowsAffected
	}

	if rowsAffected < 1 {
		return types.ErrNoRowsAffected
	}

	return nil
}

func (a *ActivityRepo) GetActivitiesHistoryByAttendanceID(attendanceID, userID int64) ([]models.Activity, error) {
	var activities []models.Activity

	rows, err := a.db.Query(`
	SELECT 
		id, 
		attendance_id, 
		user_id, 
		description, 
		created_at, 
		updated_at
	FROM activities
	WHERE user_id = ? and attendance_id = ?
	`, userID, attendanceID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetActivitiesHistoryByDate",
			"error":    err.Error(),
		}).Errorln("[AttendanceRepo] Failed to get list of activity")

		return nil, types.ErrGetActivityList
	}
	defer rows.Close()

	for rows.Next() {
		activity := models.Activity{}
		err = rows.Scan(&activity.ID, &activity.AttendanceID, &activity.UserID, &activity.Description, &activity.CreatedAt, &activity.UpdatedAt)
		if err != nil {
			return nil, err
		}

		activities = append(activities, activity)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return activities, nil
}
