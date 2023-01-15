package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
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
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprint(":", cfg.Port),
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
