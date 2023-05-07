package infrastructure

import (
	"context"
	"fmt"

	"vk_tarantool_test_task/internal/app/config"
	"vk_tarantool_test_task/internal/common"
	"vk_tarantool_test_task/internal/database"
)

type Database struct {
	db database.IDatabase
}

func New(cfg *config.Config, ctx context.Context) (*Database, error) {
	db := &Database{}
	var err error

	db.db, err = database.New(cfg, ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) InsertCredentialsHandler(ctx context.Context, userID int64, service, login, password string) error {
	// Before insert, we should check that there is already exist record for this user and service
	cred, err := db.GetCredentialsHandler(ctx, userID, service)
	if err != nil {
		return err
	}

	if len(cred) == 0 {
		// If there is no records for this service insert to database
		err = db.db.InsertCredentials(ctx, userID, service, login, password)
		if err != nil {
			return fmt.Errorf("database insert error: %w", err)
		}
	} else {
		// If there is record for this service then update record
		err = db.db.UpdateCredentials(ctx, userID, service, login, password)
		if err != nil {
			return fmt.Errorf("database insert error: %w", err)
		}
	}

	return nil
}

func (db *Database) GetCredentialsHandler(ctx context.Context, userID int64, service string) ([]*common.Credentials, error) {
	creds, err := db.db.GetCredentials(ctx, userID, service)
	if err != nil {
		return nil, fmt.Errorf("databese get error: %w", err)
	}

	return creds, nil
}

func (db *Database) DeleteCredentialsHandler(ctx context.Context, userID int64, service string) error {
	err := db.db.DeleteCredentials(ctx, userID, service)
	if err != nil {
		return fmt.Errorf("database delete error: %w", err)
	}

	return nil
}
