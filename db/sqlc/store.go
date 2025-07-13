package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store interface defines all functions to execute db queries and transactions
type Store interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	GetAccount(ctx context.Context, id int64) (Account, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
	DeleteAccount(ctx context.Context, id int64) error
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)

		return err
	})

	return result, err
}

// addMoney adds money to account balances.
// It gets the current balance and updates it with the new amount.
// To avoid deadlocks, always update accounts in order of their IDs.
func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	// Always update accounts in the same order to avoid deadlocks
	if accountID1 < accountID2 {
		account1, err = updateAccountBalance(ctx, q, accountID1, amount1)
		if err != nil {
			return
		}

		account2, err = updateAccountBalance(ctx, q, accountID2, amount2)
		if err != nil {
			return
		}
	} else {
		account2, err = updateAccountBalance(ctx, q, accountID2, amount2)
		if err != nil {
			return
		}

		account1, err = updateAccountBalance(ctx, q, accountID1, amount1)
		if err != nil {
			return
		}
	}

	return
}

// updateAccountBalance atomically updates an account's balance
func updateAccountBalance(ctx context.Context, q *Queries, accountID int64, amount int64) (Account, error) {
	account, err := q.GetAccount(ctx, accountID)
	if err != nil {
		return account, err
	}

	return q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      accountID,
		Balance: account.Balance + amount,
	})
}