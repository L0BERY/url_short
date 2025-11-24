package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"shorturl.com/internal/repository"
)

var (
	ErrTooManyAttempts = errors.New("too many attempts")
)

type ShortenerService struct {
	repo repository.Repository
}

func NewShortenerService(repo repository.Repository) *ShortenerService {
	return &ShortenerService{repo: repo}
}

func (s *ShortenerService) GenerateShortCode() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *ShortenerService) SaveURL(originalURL string) (string, error) {
	var shortCode string
	for i := 0; i < 10; i++ {
		shortCode = s.GenerateShortCode()
		if !s.repo.Exists(shortCode) {
			break
		}
		if i == 9 {
			return "", ErrTooManyAttempts
		}
	}
	err := s.repo.SaveURL(originalURL, shortCode)
	if err != nil {
		return "", err
	}
	return shortCode, nil
}

func (s *ShortenerService) GetURL(shortCode string) (string, error) {
	return s.repo.GetURL(shortCode)
}

func (s *ShortenerService) GetStats(shortCode string) (int, time.Time, error) {
	return s.repo.GetStats(shortCode)
}
