package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/app"
	internalconfig "github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/logger"
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
		return
	}

	calendar := app.New(config, logg, storage)
	fmt.Print(calendar)

	server := internalhttp.NewServer(calendar.Config, calendar.Logger)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		if err := server.Stop(); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	err := server.Start()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
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
