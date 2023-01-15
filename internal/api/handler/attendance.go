package handler

import (
	"net/http"
	"strconv"

	"github.com/scys12/go-attendance/internal/models"
	"github.com/scys12/go-attendance/internal/server"
	"github.com/scys12/go-attendance/internal/service"
	"github.com/scys12/go-attendance/internal/types"
	"github.com/sirupsen/logrus"
)

type AttendanceHandler struct {
	attendanceService service.IAttendanceService
}

func NewAttendanceHandler(svc service.IAttendanceService) AttendanceHandler {
	return AttendanceHandler{
		attendanceService: svc,
	}
}

func (a *AttendanceHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.Header.Get("X-Header-UserID"))
	if err != nil {
		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	attendance, err := a.attendanceService.CheckIn(int64(userID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"function": "CheckIn",
		}).Errorln("[AttendanceHandler] Failed to check-in user")

		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	resp := models.AttendanceResponse{
		AttendanceID:   attendance.ID,
		Message:        types.SuccessCheckIn,
		AttendanceDate: attendance.AttendanceDate.Format(types.FormatOnlyDate),
		CheckInTime:    attendance.CheckIn.Format(types.FormatTime),
	}

	server.RenderResponse(w, http.StatusOK, resp)
}

func (a *AttendanceHandler) CheckOut(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.Header.Get("X-Header-UserID"))
	if err != nil {
		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	attendance, err := a.attendanceService.CheckOut(int64(userID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"function": "CheckOut",
		}).Errorln("[AttendanceHandler] Failed to check-out user")

		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	resp := models.AttendanceResponse{
		AttendanceID:   attendance.ID,
		AttendanceDate: attendance.AttendanceDate.Format(types.FormatOnlyDate),
		Message:        types.SuccessCheckOut,
		CheckOutTime:   attendance.CheckOut.Format(types.FormatTime),
	}

	server.RenderResponse(w, http.StatusOK, resp)
}

func (a *AttendanceHandler) GetAttendanceHistory(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.Header.Get("X-Header-UserID"))
	if err != nil {
		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	attendances, err := a.attendanceService.GetAttendanceHistory(int64(userID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"function": "GetAttendanceHistory",
		}).Errorln("[AttendanceHandler] Failed to get user attendances history")

		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	resp := TransformAttendancesToAttendancesResponse(attendances)

	server.RenderResponse(w, http.StatusOK, resp)
}
