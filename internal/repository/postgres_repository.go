package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"shorturl.com/pkg/logger"
)

type PostgresRepository struct {
	db     *sql.DB
	logger logger.Logger
}

type Repository interface {
	SaveURL(originalURL, shortCode string) error
	GetURL(shortCode string) (string, error)
	Exists(shortCode string) bool
	GetStats(shortCode string) (int, time.Time, error)
	ExistsOriginalURL(originalURL string) bool
	GetShortCode(originalURL string) (string, error)
	Close() error
}

func NewPostgresRepository(connectionString string, log logger.Logger) (*PostgresRepository, error) {
	log.Info("connecting to database")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Error("failed to open database", logger.Err(err))
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		log.Error("failed to ping database", logger.Err(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTable(db); err != nil {
		log.Error("failed to create table", logger.Err(err))
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	log.Info("database connection established")
	return &PostgresRepository{db: db, logger: log}, nil
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
    	id SERIAL PRIMARY KEY,
    	short_code VARCHAR(10) UNIQUE NOT NULL,
    	original_url TEXT NOT NULL,
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	click_count INTEGER DEFAULT 0
	);
	
	CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
	CREATE INDEX IF NOT EXISTS idx_created_at ON urls(created_at);
	`

	_, err := db.Exec(query)
	return err
}

func (r *PostgresRepository) SaveURL(originalURL, shortCode string) error {
	r.logger.Debug("saving URL",
		logger.String("short_code", shortCode),
		logger.String("original_url", originalURL),
	)
	query := `INSERT INTO urls (short_code, original_url) VALUES ($1, $2)`
	_, err := r.db.Exec(query, shortCode, originalURL)
	if err != nil {
		r.logger.Error("failed to save URL",
			logger.String("short_code", shortCode),
			logger.Err(err),
		)
		return fmt.Errorf("failed to save URL: %w", err)
	}
	r.logger.Info("URL saved successfully",
		logger.String("short_code", shortCode),
	)
	return nil
}

func (r *PostgresRepository) GetURL(shortCode string) (string, error) {
	var originalURL string
	query := `SELECT original_url FROM urls WHERE short_code = $1`

	err := r.db.QueryRow(query, shortCode).Scan(&originalURL)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("URL not found for code: %s", shortCode)
	}
	if err != nil {
		return "", fmt.Errorf("database error: %s", err)
	}

	err = r.incrementClickCount(shortCode)
	if err != nil {
		fmt.Printf("Failed increment click count: %v", err)
	}

	return originalURL, nil
}

func (r *PostgresRepository) GetShortCode(originalURL string) (string, error) {
	var shortCode string
	query := `SELECT short_code FROM urls WHERE original_url = $1`

	err := r.db.QueryRow(query, originalURL).Scan(&shortCode)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("URL not found: %s", originalURL)
	}
	if err != nil {
		return "", fmt.Errorf("database error: %s", err)
	}

	return shortCode, nil
}

func (r *PostgresRepository) Exists(shortCode string) bool {
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`
	var exists bool
	err := r.db.QueryRow(query, shortCode).Scan(&exists)
	return err == nil && exists
}

func (r *PostgresRepository) ExistsOriginalURL(originalURL string) bool {
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE original_url = $1)`
	var exists bool
	err := r.db.QueryRow(query, originalURL).Scan(&exists)
	return err == nil && exists
}

func (r *PostgresRepository) incrementClickCount(shortCode string) error {
	query := `UPDATE urls SET click_count = click_count + 1 WHERE short_code = $1`
	_, err := r.db.Exec(query, shortCode)
	if err != nil {
		return fmt.Errorf("failed increment click count: %w", err)
	}
	return nil
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
