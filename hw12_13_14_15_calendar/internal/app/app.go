package app

import (
	"context"

	internalconfig "github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Config  internalconfig.Config
	Logger  logger.Logger
	Storage storage.Storage
}

func New(config internalconfig.Config, logger logger.Logger, storage storage.Storage) App {
	return App{
		config,
		logger,
		storage,
	}
}

//nolint:all
func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
