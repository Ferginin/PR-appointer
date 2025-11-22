package storage

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"PR-appointer/config"
)

func NewConnection(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	env := cfg.Env

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		env.DBUsername,
		env.DBPassword,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("Unable to parse config:", err.Error())
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		slog.Error("Unable to connect to database:", err.Error(), nil)
		panic(err)
	}

	if err = Migrate(conn); err != nil {
		slog.Error("Unable to migrate database:", err.Error(), nil)
		panic(err)
	}
	if err = DataInsert(conn); err != nil {
		slog.Error("Unable to migrate data:", err.Error(), nil)
		panic(err)
	}

	slog.Info("Connected to database")

	return conn
}
