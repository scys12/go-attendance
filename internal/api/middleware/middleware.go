package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/scys12/go-attendance/internal/models"
	"github.com/scys12/go-attendance/internal/repositories"
	"github.com/scys12/go-attendance/internal/server"
	"github.com/scys12/go-attendance/internal/types"
	"github.com/scys12/go-attendance/pkg/sessions"
	"github.com/scys12/go-attendance/pkg/tokenizer"
	"github.com/sirupsen/logrus"
)

const (
	headerEmail  = "X-Header-Email"
	headerUserID = "X-Header-UserID"
)

func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := sessions.GetSessionToken(r)
		if accessToken == "" {
			server.RenderError(w, http.StatusUnauthorized, types.ErrTokenAuthentication)
			return
		}
		accessDetail, err := tokenizer.DecodeToken(accessToken)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":    err.Error(),
				"function": "DecodeToken",
			}).Errorln("[TokenMiddleware] Failed to decode token")

			server.RenderError(w, http.StatusUnauthorized, err)
			return
		}
		userID := int(accessDetail["user_id"].(float64))
		strUserID := strconv.Itoa(userID)
		r.Header.Set(headerUserID, strUserID)

		email := accessDetail["email"].(string)
		r.Header.Set(headerEmail, email)
		next.ServeHTTP(w, r)
	})
}

func CheckInMiddleware(repo repositories.IAttendanceRepo, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(r.Header.Get("X-Header-UserID"))
		if err != nil {
			server.RenderError(w, http.StatusBadRequest, err)
			return
		}

		attendance := models.Attendance{
			AttendanceDate: time.Now(),
			UserID:         int64(userID),
		}
		checkIn, _ := repo.IsCheckIn(attendance)
		if !checkIn {
			server.RenderError(w, http.StatusBadRequest, types.ErrNotCheckIn)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CheckOutMiddleware(repo repositories.IAttendanceRepo, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(r.Header.Get("X-Header-UserID"))
		if err != nil {
			server.RenderError(w, http.StatusBadRequest, err)
			return
		}

		attendance := models.Attendance{
			AttendanceDate: time.Now(),
			UserID:         int64(userID),
		}
		checkOut, _ := repo.IsCheckOut(attendance)
		if checkOut {
			server.RenderError(w, http.StatusBadRequest, types.ErrAlreadyCheckOut)
			return
		}
		next.ServeHTTP(w, r)
	})
}
