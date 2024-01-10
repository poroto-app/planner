package rdb

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	GetDB() *sql.DB
}

func runTransaction(ctx context.Context, r Repository, f func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := r.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = f(ctx, tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
