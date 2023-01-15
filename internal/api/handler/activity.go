package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/scys12/go-attendance/internal/models"
	"github.com/scys12/go-attendance/internal/server"
	"github.com/scys12/go-attendance/internal/service"
	"github.com/scys12/go-attendance/internal/types"
	"github.com/sirupsen/logrus"
)

type ActivityHandler struct {
	activityService service.IActivityService
}

func NewActivityHandler(svc service.IActivityService) ActivityHandler {
	return ActivityHandler{
		activityService: svc,
	}
}

func (a *ActivityHandler) AddActivity(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.Header.Get("X-Header-UserID"))
	if err != nil {
		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	var req models.ActivityRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"function": "AddActivity",
		}).Info("[ActivityHandler] Unable to decode request body")

		server.RenderError(w, http.StatusBadRequest, types.ErrBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	errs := ValidateRequest(validate, err)
	if errs != nil {
		logrus.WithFields(logrus.Fields{
			"error":    errs[0],
			"function": "AddActivity",
		}).Info("[ActivityHandler] Request validation failed")

		server.RenderError(w, http.StatusBadRequest, errs[0])
		return
	}

	activity, err := a.activityService.AddActivity(req, int64(userID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"function": "AddActivity",
		}).Errorln("[ActivityHandler] Failed to insert activity")

		server.RenderError(w, http.StatusInternalServerError, err)
		return
	}

	resp := &models.ActivityResponse{
		ID:           activity.ID,
		AttendanceID: activity.AttendanceID,
		Description:  activity.Description,
	}

	server.RenderResponse(w, http.StatusCreated, resp)
}

func (a *ActivityHandler) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	activityID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := strconv.Atoi(r.Header.Get("X-Header-UserID"))
	if err != nil {
		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	var req models.ActivityRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"function": "UpdateActivity",
		}).Info("[ActivityHandler] Unable to decode request body")

		server.RenderError(w, http.StatusBadRequest, types.ErrBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	errs := ValidateRequest(validate, err)
	if errs != nil {
		logrus.WithFields(logrus.Fields{
			"error":    errs[0],
			"function": "UpdateActivity",
		}).Info("[ActivityHandler] Request validation failed")

		server.RenderError(w, http.StatusBadRequest, errs[0])
		return
	}

	activity, err := a.activityService.UpdateActivity(req, activityID, int64(userID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"function": "UpdateActivity",
		}).Errorln("[ActivityHandler] Failed to update activity")

		server.RenderError(w, http.StatusInternalServerError, err)
		return
	}

	resp := &models.ActivityResponse{
		ID:          activity.ID,
		Description: activity.Description,
	}

	server.RenderResponse(w, http.StatusOK, resp)
}

func (a *ActivityHandler) DeleteActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	activityID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := strconv.Atoi(r.Header.Get("X-Header-UserID"))
	if err != nil {
		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	err = a.activityService.DeleteActivity(activityID, int64(userID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"function": "DeleteActivity",
		}).Errorln("[ActivityHandler] Failed to delete activity")

		server.RenderError(w, http.StatusInternalServerError, err)
		return
	}

	resp := models.ActivityResponse{
		ID: activityID,
	}

	server.RenderResponse(w, http.StatusOK, resp)
}

func (a *ActivityHandler) GetActivitiesHistoryByDate(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("activity_date")
	activityDate, err := time.Parse(types.FormatOnlyDate, param)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"function": "GetActivitiesHistoryByDate",
		}).Errorln("[ActivityHandler] No query param activity_date found (" + param + ")")

		activityDate = time.Now()
	}

	userID, err := strconv.Atoi(r.Header.Get("X-Header-UserID"))
	if err != nil {
		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	activities, err := a.activityService.GetActivitiesHistoryByDate(activityDate, int64(userID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"function": "GetActivitiesHistoryByDate",
		}).Errorln("[ActivityHandler] Failed to get activities history")

		server.RenderError(w, http.StatusInternalServerError, err)
		return
	}

	resp := &models.ActivitiesResponse{
		Activities:      activities,
		TotalActivities: len(activities),
	}
	server.RenderResponse(w, http.StatusOK, resp)
}
