package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"vk_tarantool_test_task/internal/app/tgbot/repository/credentialsRepository"
	"vk_tarantool_test_task/internal/app/tgbot/service/credentialsService"

	"vk_tarantool_test_task/internal/app/config"
	"vk_tarantool_test_task/internal/app/tgbot"
)

type App struct {
	cfg *config.Config
	ctx context.Context
	db  *credentialsService.Database
	bot *tgbot.TgBot
}

func New() (*App, error) {
	app := &App{}
	var err error

	app.cfg, err = config.New()
	if err != nil {
		return nil, err
	}

	conn, err := app.connectToDB()
	if err != nil {
		return nil, err
	}

	repo, err := credentialsRepository.New(conn)
	if err != nil {
		return nil, err
	}

	app.db, err = credentialsService.New(repo)
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

func (a *App) connectToDB() (*pgxpool.Pool, error) {
	a.ctx = context.Background()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		a.cfg.Database.Host, a.cfg.Database.User, a.cfg.Database.Password, a.cfg.Database.DatabaseName, a.cfg.Database.Port)

	conn, err := pgxpool.Connect(a.ctx, dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
