package database

import (
	"context"
	"database/sql"
	"fmt"

	"aiub-companion/internal/database/db"
)

func RunInTx(ctx context.Context, conn *sql.DB, queries *db.Queries, fn func(*db.Queries) error) error {
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	qtx := queries.WithTx(tx)
	if err := fn(qtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback failed: %v, original error: %w", rbErr, err)
		}
		return err
	}
	return tx.Commit()
}
