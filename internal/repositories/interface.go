package repositories

import (
	"time"

	"github.com/scys12/go-attendance/internal/models"
)

type IAttendanceRepo interface {
	CheckIn(models.Attendance) (int64, error)
	IsCheckIn(models.Attendance) (bool, error)
	IsCheckOut(models.Attendance) (bool, error)
	CheckOut(models.Attendance) error
	GetAttendanceHistory(int64) ([]models.Attendance, error)
	GetAttendanceIDByDate(time.Time, int64) (int64, error)
}

type IUserRepo interface {
	IsEmailExist(string) (bool, error)
	InsertUser(models.User) (int64, error)
	GetUserByEmail(string) (*models.User, error)
}

type IActivityRepo interface {
	AddActivity(models.Activity) (int64, error)
	UpdateActivity(models.Activity) error
	DeleteActivity(int64, int64) error
	GetActivitiesHistoryByAttendanceID(int64, int64) ([]models.Activity, error)
}
