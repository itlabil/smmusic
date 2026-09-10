package service

import (
	"errors"

	"github.com/itlabil/smmusic/backend/internal/models"
	"github.com/itlabil/smmusic/backend/internal/repository"
)

var ErrPlaylistNotOwned = errors.New("playlist does not belong to this user")

type PlaylistService struct {
	repo *repository.PlaylistRepository
}

func NewPlaylistService(repo *repository.PlaylistRepository) *PlaylistService {
	return &PlaylistService{repo: repo}
}

func (s *PlaylistService) Create(userID int, name string) (*models.Playlist, error) {
	p := &models.Playlist{UserID: userID, Name: name}
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PlaylistService) ListByUser(userID int) ([]models.Playlist, error) {
	return s.repo.ListByUser(userID)
}

// checkOwnership ensures the playlist belongs to the requesting user before any mutation
func (s *PlaylistService) checkOwnership(playlistID, userID int) error {
	p, err := s.repo.FindByID(playlistID)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("playlist not found")
	}
	if p.UserID != userID {
		return ErrPlaylistNotOwned
	}
	return nil
}

func (s *PlaylistService) Rename(playlistID, userID int, name string) error {
	if err := s.checkOwnership(playlistID, userID); err != nil {
		return err
	}
	return s.repo.Rename(playlistID, name)
}

func (s *PlaylistService) Delete(playlistID, userID int) error {
	if err := s.checkOwnership(playlistID, userID); err != nil {
		return err
	}
	return s.repo.Delete(playlistID)
}

func (s *PlaylistService) AddSong(playlistID, userID, songID int) error {
	if err := s.checkOwnership(playlistID, userID); err != nil {
		return err
	}
	return s.repo.AddSong(playlistID, songID)
}

func (s *PlaylistService) RemoveSong(playlistID, userID, songID int) error {
	if err := s.checkOwnership(playlistID, userID); err != nil {
		return err
	}
	return s.repo.RemoveSong(playlistID, songID)
}

func (s *PlaylistService) ListSongs(playlistID, userID int) ([]models.Song, error) {
	if err := s.checkOwnership(playlistID, userID); err != nil {
		return nil, err
	}
	return s.repo.ListSongs(playlistID)
}

func (s *PlaylistService) Reorder(playlistID, userID int, songIDs []int) error {
	if err := s.checkOwnership(playlistID, userID); err != nil {
		return err
	}
	return s.repo.Reorder(playlistID, songIDs)
}