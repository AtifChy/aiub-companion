package database

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

func RunInTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("Failed to rollback transaction", "error", err)
		}
	}()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
