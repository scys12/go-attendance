package repositories

import (
	"database/sql"
	"time"

	"github.com/scys12/go-attendance/internal/models"
	"github.com/scys12/go-attendance/internal/types"
	"github.com/sirupsen/logrus"
)

type AttendanceRepo struct {
	db *sql.DB
}

func NewAttendanceRepo(db *sql.DB) IAttendanceRepo {
	return &AttendanceRepo{
		db: db,
	}
}

func (a *AttendanceRepo) CheckIn(attendance models.Attendance) (int64, error) {
	var id int64
	res, err := a.db.Exec(`
		INSERT INTO attendances(
			check_in,
			attendance_date,
			user_id
		) VALUES (
			?,
			?,
			?
		)`,
		attendance.CheckIn,
		attendance.AttendanceDate,
		attendance.UserID,
	)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "CheckIn",
			"error":    err.Error(),
		}).Errorln("[AttendanceRepo] Failed to insert attendance")
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "CheckIn",
			"error":    err.Error(),
		}).Errorln("[AttendanceRepo] Failed to get ID")

		return 0, types.ErrFailedGetLastIDDB
	}

	return id, nil
}

func (a *AttendanceRepo) IsCheckIn(attendance models.Attendance) (bool, error) {
	isExist := false
	attendanceDate := attendance.AttendanceDate.Format(types.FormatOnlyDate)

	if err := a.db.QueryRow(`
		SELECT check_in 
		FROM attendances 
		WHERE attendance_date = ? and user_id = ?
		`,
		attendanceDate,
		attendance.UserID,
	).Scan(&attendance.CheckIn); err != nil {
		return false, err
	}

	if !attendance.CheckIn.IsZero() {
		isExist = true
	}
	return isExist, nil
}

func (a *AttendanceRepo) IsCheckOut(attendance models.Attendance) (bool, error) {
	isExist := false
	attendanceDate := attendance.AttendanceDate.Format(types.FormatOnlyDate)

	if err := a.db.QueryRow(`
		SELECT check_out 
		FROM attendances 
		WHERE attendance_date = ? and user_id = ?
		`,
		attendanceDate,
		attendance.UserID,
	).Scan(&attendance.CheckOut); err != nil {
		return false, err
	}

	if !attendance.CheckOut.IsZero() {
		isExist = true
	}
	return isExist, nil
}

func (a *AttendanceRepo) CheckOut(attendance models.Attendance) error {
	res, err := a.db.Exec(`
		UPDATE attendances
		SET check_out = ?
		WHERE id = ? and user_id = ?`,
		attendance.CheckOut,
		attendance.ID,
		attendance.UserID,
	)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "CheckOut",
			"error":    err.Error(),
		}).Errorln("[AttendanceRepo] Failed to update attendance")
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "UpdateCake",
			"error":    err.Error(),
		}).Errorln("[CakeRepo] Failed to get rows affected")

		return types.ErrFailedRowsAffected
	}

	if rowsAffected < 1 {
		return types.ErrNoRowsAffected
	}

	return nil
}

func (a *AttendanceRepo) GetAttendanceHistory(userID int64) ([]models.Attendance, error) {
	var attendances []models.Attendance
	var checkoutTime sql.NullTime

	rows, err := a.db.Query(`
	SELECT 
		id, 
		check_in, 
		check_out, 
		attendance_date, 
		user_id, 
		created_at, 
		updated_at
	FROM attendances
	WHERE user_id = ?
	`, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetAttendanceHistory",
			"error":    err.Error(),
		}).Errorln("[AttendanceRepo] Failed to get list of attendance")

		return nil, types.ErrGetAttendanceList
	}
	defer rows.Close()

	for rows.Next() {
		attendance := models.Attendance{}
		err = rows.Scan(&attendance.ID, &attendance.CheckIn, &checkoutTime, &attendance.AttendanceDate, &attendance.UserID, &attendance.CreatedAt, &attendance.UpdatedAt)
		if err != nil {
			return nil, err
		}
		attendance.CheckOut = checkoutTime.Time
		attendances = append(attendances, attendance)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return attendances, nil
}

func (a *AttendanceRepo) GetAttendanceIDByDate(attDate time.Time, userID int64) (int64, error) {
	var attendanceID int64
	attendanceDate := attDate.Format(types.FormatOnlyDate)

	if err := a.db.QueryRow(`
		SELECT id
		FROM attendances 
		WHERE attendance_date = ? and user_id = ?
		`,
		attendanceDate,
		userID,
	).Scan(&attendanceID); err != nil {
		return 0, err
	}

	return attendanceID, nil
}
