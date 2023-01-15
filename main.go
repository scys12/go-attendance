package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/scys12/go-attendance/config"
	"github.com/scys12/go-attendance/database"
	"github.com/scys12/go-attendance/internal/api/handler"
	"github.com/scys12/go-attendance/internal/api/middleware"
	"github.com/scys12/go-attendance/internal/repositories"
	"github.com/scys12/go-attendance/internal/server"
	"github.com/scys12/go-attendance/internal/service"
	"github.com/scys12/go-attendance/pkg/sessions"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	dbConfig, err := config.InitDBConfig()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"config": dbConfig,
		}).Error("[Config] Failed to load server config")
	}
	servConf, err := config.InitServerConfig()
	if err != nil {
		logrus.Info("[Config] Failed to load server config")
	}

	db := database.GetDatabaseConnection(dbConfig)

	userRepo := repositories.NewUserRepo(db)
	attendanceRepo := repositories.NewAttendanceRepo(db)
	activityRepo := repositories.NewActivityRepo(db)

	userService := service.NewUserService(userRepo)
	attendanceService := service.NewAttendanceService(attendanceRepo)
	acitivtyService := service.NewActivityService(activityRepo, attendanceRepo)

	userHandler := handler.NewUserHandler(userService)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)
	activityHandler := handler.NewActivityHandler(acitivtyService)

	router := mux.NewRouter()
	router.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/register", userHandler.Register).Methods(http.MethodPost)
	router.Path("/logout").Methods(http.MethodPost).Handler(middleware.TokenMiddleware(http.HandlerFunc(userHandler.Logout)))

	attendanceRtr := router.PathPrefix("/attendance").Subrouter()
	attendanceRtr.Use(middleware.TokenMiddleware)
	attendanceRtr.HandleFunc("", attendanceHandler.GetAttendanceHistory).Methods(http.MethodGet)
	attendanceRtr.HandleFunc("/checkin", attendanceHandler.CheckIn).Methods(http.MethodPost)
	attendanceRtr.Path("/checkout").Methods(http.MethodPost).Handler(middleware.CheckOutMiddleware(attendanceRepo, middleware.CheckInMiddleware(attendanceRepo, http.HandlerFunc(attendanceHandler.CheckOut))))

	activityRtr := router.PathPrefix("/activity").Subrouter()
	activityRtr.Use(middleware.TokenMiddleware)
	activityRtr.Path("").Methods(http.MethodPost).Handler(middleware.CheckInMiddleware(attendanceRepo, http.HandlerFunc(activityHandler.AddActivity)))
	activityRtr.HandleFunc("", activityHandler.GetActivitiesHistoryByDate).Methods(http.MethodGet)
	activityRtr.Path("/{id}").Methods(http.MethodPatch).Handler(middleware.CheckInMiddleware(attendanceRepo, http.HandlerFunc(activityHandler.UpdateActivity)))
	activityRtr.Path("/{id}").Methods(http.MethodDelete).Handler(middleware.CheckInMiddleware(attendanceRepo, http.HandlerFunc(activityHandler.DeleteActivity)))

	sessionManager := sessions.GetSessionInstance()
	wrappedMux := sessionManager.LoadAndSave(router)

	serverConfig := server.Config{
		WriteTimeout: time.Duration(servConf.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(servConf.ReadTimeout) * time.Second,
		Port:         servConf.Port,
	}
	server.Serve(serverConfig, wrappedMux)
}
