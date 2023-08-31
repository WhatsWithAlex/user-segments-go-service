package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/WhatsWithAlex/user-segments-go-service/internal/env"
	"github.com/WhatsWithAlex/user-segments-go-service/internal/postgresdb"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupDB(env *env.DBEnv) (db *postgresdb.Store, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		env.User,
		env.Password,
		env.Host,
		env.Port,
		env.Name,
	)

	err = runMigrate(connString)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return
	}

	dbConn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return
	}
	db = postgresdb.NewStore(dbConn)
	return
}

func runMigrate(connString string) error {
	m, err := migrate.New(
		"file://./sql/migrations",
		connString,
	)
	if err != nil {
		return err
	}
	err = m.Up()
	return err
}
