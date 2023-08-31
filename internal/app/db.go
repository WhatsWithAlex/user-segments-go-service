package app

import (
	"context"
	"fmt"
	"time"

	"github.com/WhatsWithAlex/user-segments-go-service/internal/env"
	"github.com/WhatsWithAlex/user-segments-go-service/internal/postgresdb"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupDB(env *env.DBEnv) (db *postgresdb.Store, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		env.User,
		env.Password,
		env.Addr,
		env.Port,
		env.Name,
	)
	dbConn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return
	}
	db = postgresdb.NewStore(dbConn)
	return
}
