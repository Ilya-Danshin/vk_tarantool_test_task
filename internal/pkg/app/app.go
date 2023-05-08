package app

import (
	"context"

	"vk_tarantool_test_task/internal/app/config"
	"vk_tarantool_test_task/internal/app/tgbot"
	"vk_tarantool_test_task/internal/infrastructure"
)

type App struct {
	cfg *config.Config
	ctx context.Context
	db  *infrastructure.Database
	bot *tgbot.TgBot
}

func New() (*App, error) {
	app := &App{}
	var err error

	app.cfg, err = config.New()
	if err != nil {
		return nil, err
	}

	app.ctx = context.Background()
	app.db, err = infrastructure.New(app.cfg, app.ctx)
	if err != nil {
		return nil, err
	}

	app.bot, err = tgbot.New(app.cfg, app.db)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run() error {
	err := a.bot.Run(a.ctx)
	if err != nil {
		return err
	}

	return nil
}
