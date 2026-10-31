package main

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrInvalidAmount     = errors.New("invalid transfer amount")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrAccountNotFound   = errors.New("account not found")
)

func TransferMoney(
	ctx context.Context,
	db *sql.DB,
	fromID int64,
	toID int64,
	amount int64,
) error {

	// TODO 1: Validate the amount
	if amount <= 0 {
		return ErrInvalidAmount
	}

	// TODO 2: Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// TODO 3: Decrease sender's balance
	res, err := tx.ExecContext(ctx, "UPDATE accounts SET balance_paise = balance_paise - ? WHERE id = ? AND balance_paise >= ?", amount, fromID, amount)
	if err != nil {
		return err
	}

	// TODO 4: Check RowsAffected()
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrAccountNotFound
	}

	// TODO 5: Increase receiver's balance
	res2, err := tx.ExecContext(ctx, "UPDATE accounts SET balance_paise = balance_paise + ? WHERE id = ?", amount, toID)
	if err != nil {
		return err
	}

	rowsAffected2, err := res2.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected2 == 0 {
		return ErrAccountNotFound
	}

	// TODO 6: Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
