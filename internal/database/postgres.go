package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"vk_tarantool_test_task/internal/app/config"
	"vk_tarantool_test_task/internal/common"
)

type Postgres struct {
	conn *pgxpool.Pool
}

func New(cfg *config.Config, ctx context.Context) (*Postgres, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DatabaseName, cfg.Database.Port)

	conn, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	db := Postgres{}
	db.conn = conn

	return &db, nil
}

const insertIntoCredentialsTable = `INSERT INTO credentials(user_id, service, login, password) 
									VALUES ($1, $2, $3, $4);`

func (db *Postgres) InsertCredentials(ctx context.Context, userID int64, service, login, password string) error {
	row, err := db.conn.Query(ctx, insertIntoCredentialsTable, userID, service, login, password)
	if err != nil {
		return err
	}
	defer row.Close()

	return nil
}

const updateCredentialsTable = `UPDATE credentials
								SET login=$1, password=$2
								WHERE user_id=$3 AND service = $4`

func (db *Postgres) UpdateCredentials(ctx context.Context, userID int64, service, login, password string) error {
	row, err := db.conn.Query(ctx, updateCredentialsTable, login, password, userID, service)
	if err != nil {
		return err
	}
	defer row.Close()

	return nil
}

const getFromCredentialsTable = `SELECT user_id, service, login, password 
									FROM credentials
									WHERE user_id=$1 AND service=$2;`

func (db *Postgres) GetCredentials(ctx context.Context, userID int64, service string) ([]*common.Credentials, error) {
	rows, err := db.conn.Query(ctx, getFromCredentialsTable, userID, service)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var results []*common.Credentials

	for rows.Next() {
		var result common.Credentials

		err = rows.Scan(
			&result.UserID,
			&result.Service,
			&result.Login,
			&result.Password,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &result)
	}

	return results, nil
}

const deleteFromCredentialsTable = `DELETE FROM credentials
									WHERE user_id=$1 AND service=$2;`

func (db *Postgres) DeleteCredentials(ctx context.Context, userID int64, service string) error {
	row, err := db.conn.Query(ctx, deleteFromCredentialsTable, userID, service)
	if err != nil {
		return err
	}
	defer row.Close()

	return nil
}
