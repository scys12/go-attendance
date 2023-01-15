package service

import (
	"time"

	"github.com/scys12/go-attendance/internal/models"
	"github.com/scys12/go-attendance/internal/repositories"
	"github.com/scys12/go-attendance/internal/types"
	"github.com/sirupsen/logrus"
)

type AttendanceService struct {
	attendanceRepo repositories.IAttendanceRepo
}

func NewAttendanceService(repo repositories.IAttendanceRepo) IAttendanceService {
	return &AttendanceService{
		attendanceRepo: repo,
	}
}

func (a *AttendanceService) CheckIn(userID int64) (*models.Attendance, error) {
	now := time.Now()
	attendance := models.Attendance{
		CheckIn:        now,
		UserID:         userID,
		AttendanceDate: now,
	}

	checkIn, _ := a.attendanceRepo.IsCheckIn(attendance)
	if checkIn {
		return nil, types.ErrAlreadyCheckIn
	}

	attendanceID, err := a.attendanceRepo.CheckIn(attendance)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "CheckIn",
			"error":    err.Error(),
		}).Errorln("[AttendanceService] Failed to check in user")

		return nil, err
	}

	attendance.ID = attendanceID

	return &attendance, nil
}

func (a *AttendanceService) CheckOut(userID int64) (*models.Attendance, error) {
	now := time.Now()
	attendanceID, err := a.attendanceRepo.GetAttendanceIDByDate(now, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "CheckOut",
			"error":    err.Error(),
		}).Errorln("[AttendanceService] Failed to get attendance id by date")

		return nil, err
	}

	attendance := models.Attendance{
		CheckOut: now,
		UserID:   userID,
		ID:       attendanceID,
	}
	err = a.attendanceRepo.CheckOut(attendance)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "CheckOut",
			"error":    err.Error(),
		}).Errorln("[AttendanceService] Failed to check out user")

		return nil, err
	}

	attendance.AttendanceDate = now

	return &attendance, nil
}

func (a *AttendanceService) GetAttendanceHistory(userID int64) ([]models.Attendance, error) {
	attendances, err := a.attendanceRepo.GetAttendanceHistory(userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetAttendanceHistory",
			"error":    err.Error(),
		}).Errorln("[AttendanceService] Failed to get list of attendance")

		return nil, err
	}

	return attendances, nil
}
