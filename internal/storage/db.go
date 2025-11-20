package storage

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed init.sql
var initSQL string

func CheckAndMigrate(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), initSQL)
	if err != nil {
		return fmt.Errorf("init sql failed: %v", err)
	}

	var count int
	err = db.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("get users count failed: %v", err)
	}

	if count == 0 {
		if err = InsertAdminUser(context.Background(), db); err != nil {
			return fmt.Errorf("insert admin user failed: %v", err)
		}
	}

	return nil
}

func InsertAdminUser(ctx context.Context, db *pgxpool.Pool) error {
	// hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	_, err := db.Exec(
		context.Background(),
		`INSERT INTO users (id, username)
		VALUES (DEFAULT, $1)
		RETURNING id
		`,
		"admin",
	)

	return err
}
