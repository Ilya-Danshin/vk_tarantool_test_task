package app

import (
	"vk_tarantool_test_task/internal/app/config"
	"vk_tarantool_test_task/internal/app/tgbot"
)

type App struct {
	cfg *config.Config
	bot *tgbot.TgBot
}

func New() (*App, error) {
	app := &App{}
	var err error

	app.cfg, err = config.New()
	if err != nil {
		return nil, err
	}

	app.bot, err = tgbot.New(app.cfg)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run() error {
	err := a.bot.Run()
	if err != nil {
		return err
	}

	return nil
}
