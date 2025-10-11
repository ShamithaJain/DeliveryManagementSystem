package internal

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func NewTestStore() (*Store, error) {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        return nil, err
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT NOT NULL UNIQUE,
            password_hash TEXT NOT NULL,
            role TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );
        CREATE TABLE IF NOT EXISTS orders (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            customer_id INTEGER NOT NULL,
            items TEXT NOT NULL, -- store JSON as TEXT
            status TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            cancelled_at DATETIME
    );


    `)
    if err != nil {
        return nil, err
    }

    return &Store{
        DB:    db,
        Cache: make(map[int]*Order),
    }, nil
}
