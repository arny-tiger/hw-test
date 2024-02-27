package internalgrpc

import (
	"net"
	"strconv"

	pbEvents "github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/api/grpc/gen/event"
	eventsHandlers "github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/api/grpc/handler/event"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	logger     logger.Logger
	config     config.GRPCConf
}

func NewServer(config config.Config, logger logger.Logger, storage storage.Storage) Server {
	grpcsServ := grpc.NewServer(grpc.UnaryInterceptor(getLoggingInterceptor(logger)))
	reflection.Register(grpcsServ)
	s := Server{
		grpcServer: grpcsServ,
		logger:     logger,
		config:     config.GRPC,
	}
	eventsHandler := &eventsHandlers.Handler{Storage: storage, Logger: logger}
	pbEvents.RegisterEventServiceServer(s.grpcServer, eventsHandler)
	return s
}

func (s *Server) Start() error {
	port := ":" + strconv.Itoa(s.config.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	err = s.grpcServer.Serve(lis)
	if err != nil {
		return err
	}
	s.logger.Info("gRPC server listening at " + lis.Addr().String())
	return nil
}

func (s *Server) Stop() {
	s.grpcServer.Stop()
}
