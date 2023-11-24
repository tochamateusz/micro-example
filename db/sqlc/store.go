package simple_bank

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func ProvideStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

var txKey = struct{}{}

type QueryFnc = func(*Queries) error

func (store *Store) execTx(ctx context.Context, fn QueryFnc) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:%v, rb err: %v ", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer      Transfer `json:"transfer"`
	FromAccountID Account  `json:"from_account_id"`
	ToAccountID   Account  `json:"to_account_id"`
	FromEntry     Entry    `json:"from_entry"`
	ToEntry       Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) (err error) {

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		fromAccount, err := q.GetAccountByIDForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		result.FromAccountID, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: fromAccount.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		toAccount, err := q.GetAccountByIDForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		result.ToAccountID, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: toAccount.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
