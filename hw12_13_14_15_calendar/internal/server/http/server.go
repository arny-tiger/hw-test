package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	httpServer *http.Server
	logger     logger.Logger
}

func NewServer(config config.Config, logger logger.Logger) Server {
	httpServer := &http.Server{
		Addr:         config.Host.Host + ":" + strconv.Itoa(config.Host.Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the home page!")
	})
	httpServer.Handler = latencyMiddleware(loggingMiddleware(mux, logger))

	return Server{
		httpServer,
		logger,
	}
}

func (s *Server) Start() error {
	err := s.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		s.logger.Error("Server Shutdown Failed: " + err.Error())
		return err
	}
	return nil
}
