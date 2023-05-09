package repository

import (
	"context"

	"vk_tarantool_test_task/internal/common"
)

type CredentialsRepository interface {
	Insert(ctx context.Context, userID int64, service, login, password string) error
	Select(ctx context.Context, userID int64, service string) ([]*common.Credentials, error)
	Delete(ctx context.Context, userID int64, service string) error
}
