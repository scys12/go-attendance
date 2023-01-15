package service

import (
	"time"

	"github.com/scys12/go-attendance/internal/models"
	"github.com/scys12/go-attendance/internal/repositories"
	"github.com/sirupsen/logrus"
)

type ActivityService struct {
	activityRepo   repositories.IActivityRepo
	attendanceRepo repositories.IAttendanceRepo
}

func NewActivityService(activityRepo repositories.IActivityRepo, attendanceRepo repositories.IAttendanceRepo) IActivityService {
	return &ActivityService{
		activityRepo:   activityRepo,
		attendanceRepo: attendanceRepo,
	}
}

func (a *ActivityService) AddActivity(req models.ActivityRequest, userID int64) (*models.Activity, error) {
	attendanceID, err := a.attendanceRepo.GetAttendanceIDByDate(time.Now(), userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "AddActivity",
			"error":    err.Error(),
		}).Errorln("[ActivityService] Failed to get attendance id by date")

		return nil, err
	}

	activity := models.Activity{
		Description:  req.Description,
		AttendanceID: attendanceID,
		UserID:       userID,
	}

	id, err := a.activityRepo.AddActivity(activity)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "AddActivity",
			"error":    err.Error(),
		}).Errorln("[ActivityService] Failed to add activity")

		return nil, err
	}

	activity.ID = id

	return &activity, nil
}

func (a *ActivityService) UpdateActivity(req models.ActivityRequest, activityID, userID int64) (*models.Activity, error) {
	activity := models.Activity{
		ID:          activityID,
		Description: req.Description,
		UserID:      userID,
	}

	err := a.activityRepo.UpdateActivity(activity)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "UpdateActivity",
			"error":    err.Error(),
		}).Errorln("[ActivityService] Failed to update activity")

		return nil, err
	}

	return &activity, nil
}

func (a *ActivityService) DeleteActivity(activityID, userID int64) error {
	err := a.activityRepo.DeleteActivity(activityID, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "DeleteActivity",
			"error":    err.Error(),
		}).Errorln("[ActivityService] Failed to delete activity")

		return err
	}

	return nil
}

func (a *ActivityService) GetActivitiesHistoryByDate(activityDate time.Time, userID int64) ([]models.Activity, error) {
	attendanceID, err := a.attendanceRepo.GetAttendanceIDByDate(activityDate, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetActivitiesHistoryByDate",
			"error":    err.Error(),
		}).Errorln("[ActivityService] Failed to get attendance id by date")

		return []models.Activity{}, nil
	}

	activities, err := a.activityRepo.GetActivitiesHistoryByAttendanceID(attendanceID, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetActivitiesHistoryByDate",
			"error":    err.Error(),
		}).Errorln("[ActivityService] Failed to get activities id by attendance id")

		return nil, err
	}

	return activities, nil
}
