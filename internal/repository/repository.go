package repository

import (
	"time"
)

// var (
// 	ErrNotFound = errors.New("URL not found")
// )

type Repository interface {
	SaveURL(originalURL, shortCode string) error
	GetURL(shortCode string) (string, error)
	Exist(shortCode string) bool
	GetStats(shortCode string) (int, time.Time, error)
	Close() error
}

//	ТУТ БД ПОДРУБАЕМ
