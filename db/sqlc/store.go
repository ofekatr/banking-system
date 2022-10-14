package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// Store encapsulates all the database operations.
type Store struct {
	*Queries
	db *sql.DB
}

// Creats a new store.
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a function within a database transaction.
func (store *Store) execTx(ctx context.Context, txFunc func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin a transaction")
	}

	q := New(tx)
	if txErr := txFunc(q); txErr != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to execute transaction: %v, failed to rollback transaction: %v", txErr, rbErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (*TransferTxResult, error) {
	result := TransferTxResult{}

	err := store.execTx(ctx, func(q *Queries) (err error) {
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return errors.Wrap(err, "failed to create transfer")
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create from entry")
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create to entry")
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.AddBalanceAmount(ctx, AddBalanceAmountParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return errors.Wrap(err, "failed to add amount to account")
			}
			result.ToAccount, err = q.AddBalanceAmount(ctx, AddBalanceAmountParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return errors.Wrap(err, "failed to add amount to account")
			}
		} else {
			result.ToAccount, err = q.AddBalanceAmount(ctx, AddBalanceAmountParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return errors.Wrap(err, "failed to add amount to account")
			}

			result.FromAccount, err = q.AddBalanceAmount(ctx, AddBalanceAmountParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return errors.Wrap(err, "failed to add amount to account")
			}
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute transfer transaction")
	}

	return &result, nil
}
