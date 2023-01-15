package service

import (
	"time"

	"github.com/scys12/go-attendance/internal/models"
)

type IAttendanceService interface {
	CheckIn(int64) (*models.Attendance, error)
	CheckOut(int64) (*models.Attendance, error)
	GetAttendanceHistory(int64) ([]models.Attendance, error)
}

type IUserService interface {
	RegisterUser(models.RegisterRequest) (*models.User, error)
	LoginUser(models.LoginRequest) (*models.User, error)
}

type IActivityService interface {
	AddActivity(models.ActivityRequest, int64) (*models.Activity, error)
	UpdateActivity(models.ActivityRequest, int64, int64) (*models.Activity, error)
	DeleteActivity(int64, int64) error
	GetActivitiesHistoryByDate(time.Time, int64) ([]models.Activity, error)
}
