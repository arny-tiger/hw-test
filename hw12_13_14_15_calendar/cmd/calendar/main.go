package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/app"
	internalconfig "github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/server/http"
	iStorage "github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(
		&configFile,
		"config",
		"./configs/config.yaml",
		"Path to configuration file",
	)
}

func main() {
	flag.Parse()
	config := internalconfig.NewConfig(configFile)
	logg := logger.New(config.Logger.Level)
	storage, storageErr := createStorage(config)
	if storageErr != nil {
		logg.Error(storageErr.Error())
		os.Exit(1)
	}

	calendar := app.New(config, logg, storage)
	httpServer := internalhttp.NewServer(calendar.Config, calendar.Logger, calendar.Storage)
	grpcServer := internalgrpc.NewServer(calendar.Config, calendar.Logger, calendar.Storage)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	startHTTPServer(cancel, httpServer, logg)
	logg.Info("HTTP server listening at " + config.Host.Host + ":" + strconv.Itoa(config.Host.Port))

	startGRPCServer(cancel, grpcServer, logg)
	logg.Info("gRPC server listening at " + config.GRPC.Host + ":" + strconv.Itoa(config.GRPC.Port))

	listenForCancel(ctx, logg, httpServer, grpcServer)
}

func createStorage(config internalconfig.Config) (iStorage.Storage, error) {
	if config.DB.Type == iStorage.SQLStorageType {
		sqlStorage, storageErr := sqlstorage.New(config)
		return sqlStorage, storageErr
	}
	if config.DB.Type == iStorage.MemoryStorageType {
		memoryStorage := memorystorage.New()
		return memoryStorage, nil
	}
	return nil, fmt.Errorf("wrong storage type")
}

func startHTTPServer(cancel context.CancelFunc, server internalhttp.Server, logg logger.Logger) {
	go func() {
		err := server.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}()
}

func startGRPCServer(cancel context.CancelFunc, server internalgrpc.Server, logg logger.Logger) {
	go func() {
		err := server.Start()
		if err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
	}()
}

func listenForCancel(ctx context.Context, logg logger.Logger, http internalhttp.Server, grpc internalgrpc.Server) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := http.Stop(); err != nil {
			logg.Error("failed to stop http HTTP Server: " + err.Error())
		}
		grpc.Stop()
		logg.Info("Calendar has been stopped")
	}()
	wg.Wait()
}
