package main

import (
	"context"
	"database/sql"
	"log"
)

func createTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS accounts (
            id INTEGER PRIMARY KEY,
            balance_paise INTEGER NOT NULL
        )
    `

	_, err := db.ExecContext(context.Background(), query)
	if err != nil {
		return err
	}

	query = `
        INSERT OR IGNORE INTO accounts
        (id, balance_paise)
        VALUES
        (101, 100000),
        (102, 50000)
    `

	_, err = db.ExecContext(context.Background(), query)
	return err
}

func main() {
	db, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createTable(db); err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, balance_paise FROM accounts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, balance int64

		if err := rows.Scan(&id, &balance); err != nil {
			log.Fatal(err)
		}

		log.Printf("Account: %d, Balance: %d paise", id, balance)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	TransferMoney(context.Background(), db, 101, 102, 10000)

	log.Println("Database ready!")
}
