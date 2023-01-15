package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	Port         int
}

func Serve(cfg Config, router http.Handler) {
	var port string
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	} else {
		port = strconv.Itoa(cfg.Port)
	}
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprint(":", port),
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.WithFields(logrus.Fields{
				"port":     cfg.Port,
				"error":    err.Error(),
				"function": "Serve",
			}).Fatal("[Server] Unable to listen and serve")
		}
	}()
	logrus.Info("[Server] HTTP server is running at port ", cfg.Port)

	s := make(chan os.Signal, 1)

	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-s

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Error("[Server] error on shutting down HTTP Server, err: ", err.Error())
	}
}
