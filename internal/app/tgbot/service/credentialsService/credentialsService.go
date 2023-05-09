package credentialsService

import (
	"context"
	"fmt"

	"vk_tarantool_test_task/internal/app/tgbot/repository"
	"vk_tarantool_test_task/internal/common"
)

type Database struct {
	repo repository.CredentialsRepository
}

func New(repo repository.CredentialsRepository) (*Database, error) {
	db := Database{
		repo: repo,
	}

	return &db, nil
}

func (db *Database) InsertCredentials(ctx context.Context, userID int64, service, login, password string) error {
	// Before insert, we should check that there is already exist record for this user and service
	cred, err := db.GetCredentials(ctx, userID, service)
	if err != nil {
		return err
	}

	if len(cred) == 0 {
		// If there is no records for this service insert to database
		err = db.repo.Insert(ctx, userID, service, login, password)
		if err != nil {
			return fmt.Errorf("database insert error: %w", err)
		}
	} else {
		// If there is record for this service then update record
		err = db.repo.Update(ctx, userID, service, login, password)
		if err != nil {
			return fmt.Errorf("database insert error: %w", err)
		}
	}

	return nil
}

func (db *Database) GetCredentials(ctx context.Context, userID int64, service string) ([]*common.Credentials, error) {
	creds, err := db.repo.Select(ctx, userID, service)
	if err != nil {
		return nil, fmt.Errorf("databese get error: %w", err)
	}

	return creds, nil
}

func (db *Database) DeleteCredentials(ctx context.Context, userID int64, service string) error {
	err := db.repo.Delete(ctx, userID, service)
	if err != nil {
		return fmt.Errorf("database delete error: %w", err)
	}

	return nil
}
