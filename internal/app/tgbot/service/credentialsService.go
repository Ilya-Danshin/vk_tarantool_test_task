package service

import (
	"context"

	"vk_tarantool_test_task/internal/common"
)

type CredentialsService interface {
	InsertCredentials(ctx context.Context, userID int64, service, login, password string) error
	GetCredentials(ctx context.Context, userID int64, service string) ([]*common.Credentials, error)
	DeleteCredentials(ctx context.Context, userID int64, service string) error
}
