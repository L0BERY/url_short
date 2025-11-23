package service

import "shorturl.com/internal/repository"

type ShortenerService struct {
	repo repository.Repository
}

func NewShortenerService(repo repository.Repository) *ShortenerService {
	return &ShortenerService{repo: repo}
}

//	РАБОТА С ССЫЛКАМИ
