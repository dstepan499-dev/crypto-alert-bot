package storage

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

type Alert struct {
	ID          int64
	ChatID      int64
	Coin        string
	TargetPrice float64
}

type Storage struct {
	db *sql.DB
}

func NewStorage(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("Error while opening DB: %w", err)
	}

	// Create alerts table
	query := `CREATE TABLE IF NOT EXISTS alerts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	chat_id INTEGER NOT NULL,
	coin TEXT NOT NULL,
	target_price REAL NOT NULL
	);`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("Table creating error: %w", err)
	}

	return &Storage{db: db}, nil
}

// AddAlert saves new user's submition
func (s Storage) AddAlert(chatID int64, coin string, targetPrice float64) error {
	coin = strings.ToLower(strings.TrimSpace(coin))
	query := `INSERT INTO alerts (chat_id, coin, target_price) VALUES (?, ?, ?)`
	_, err := s.db.Exec(query, chatID, coin, targetPrice)
	if err != nil {
		return fmt.Errorf("Alert insertion error: %w", err)
	}
	return nil
}

// GetAlerts returns all active alerts list
func (s *Storage) GetAlerts() ([]Alert, error) {
	query := `SELECT id, chat_id, coin, target_price FROM alerts`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Selection error: %w", err)
	}
	defer rows.Close()

	var alerts []Alert
	for rows.Next() {
		var a Alert
		if err := rows.Scan(&a.ID, &a.ChatID, &a.Coin, &a.TargetPrice); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}

	return alerts, nil
}

// DeleteAlert deletes worked alerts
func (s *Storage) DeleteAlert(id int64) error {
	query := `DELETE FROM alerts WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}
