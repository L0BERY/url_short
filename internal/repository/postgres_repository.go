package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

type Repository interface {
	SaveURL(originalURL, shortCode string) error
	GetURL(shortCode string) (string, error)
	Exists(shortCode string) bool
	GetStats(shortCode string) (int, time.Time, error)
	Close() error
}

func NewPostgresRepository(connectionString string) (*PostgresRepository, error) {
	db, err1 := sql.Open("postgres", connectionString)
	if err1 != nil {
		return nil, fmt.Errorf("failed to open database: %w", err1)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// sqlBytes, err2 := os.ReadFile()
	// if err2 != nil {
	// 	return nil, fmt.Errorf("failed to read file: %w", err2)
	// }
	// sqlQuery := string(sqlBytes)
	// _, err := db.Exec(sqlQuery)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create database: %w", err)
	// }

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) SaveURL(originalURL, shortCode string) error {
	query := `INSERT INTO urls (short_code, original_url) VALUES ($1, $2)`
	_, err := r.db.Exec(query, shortCode, originalURL)
	if err != nil {
		return fmt.Errorf("failed to save URL: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetURL(shortCode string) (string, error) {
	var originalURL string
	query := `SELECT original_url FROM urls WHERE short_code = $1`

	err := r.db.QueryRow(query, shortCode).Scan(&originalURL)
	if err != nil {
		return "", err
	}

	go r.incrementClickCount(shortCode)

	return originalURL, err
}

func (r *PostgresRepository) Exists(shortCode string) bool {
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`
	var exists bool
	err := r.db.QueryRow(query, shortCode).Scan(&exists)
	return err == nil && exists
}

func (r *PostgresRepository) incrementClickCount(shortCode string) {
	query := `UPDATE urls SET click_count = click_count + 1 WHERE short_code = $1`
	r.db.Exec(query, shortCode)
}

func (r *PostgresRepository) GetStats(shortCode string) (int, time.Time, error) {
	var clickCount int
	var createdAt time.Time
	query := `SELECT click_count, created_at FROM urls WHERE short_code = $1`
	err := r.db.QueryRow(query, shortCode).Scan(&clickCount, &createdAt)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to get ststs: %w", err)
	}
	return clickCount, createdAt, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
