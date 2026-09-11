package service

import (
	"github.com/itlabil/smmusic/backend/internal/models"
	"github.com/itlabil/smmusic/backend/internal/repository"
)

type InteractionService struct {
	repo *repository.InteractionRepository
}

func NewInteractionService(repo *repository.InteractionRepository) *InteractionService {
	return &InteractionService{repo: repo}
}

func (s *InteractionService) LikeSong(userID, songID int) error {
	return s.repo.LikeSong(userID, songID)
}

func (s *InteractionService) UnlikeSong(userID, songID int) error {
	return s.repo.UnlikeSong(userID, songID)
}

func (s *InteractionService) ListLikedSongs(userID int) ([]models.Song, error) {
	return s.repo.ListLikedSongs(userID)
}

func (s *InteractionService) RecordPlay(userID, songID int) error {
	return s.repo.RecordPlay(userID, songID)
}

func (s *InteractionService) ListRecentlyPlayed(userID int) ([]models.Song, error) {
	return s.repo.ListRecentlyPlayed(userID, 20) // last 20 unique songs
}

func (s *InteractionService) GetLikedSongsSummary(userID int) (*repository.LikedSongsSummary, error) {
	return s.repo.GetLikedSongsSummary(userID)
}