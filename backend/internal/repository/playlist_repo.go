package repository

import (
	"database/sql"
	"errors"

	"github.com/itlabil/smmusic/backend/internal/models"
)

type PlaylistRepository struct {
	db *sql.DB
}

func NewPlaylistRepository(db *sql.DB) *PlaylistRepository {
	return &PlaylistRepository{db: db}
}

func (r *PlaylistRepository) Create(p *models.Playlist) error {
	query := `
		INSERT INTO playlists (user_id, name)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, p.UserID, p.Name).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PlaylistRepository) FindByID(id int) (*models.Playlist, error) {
	query := `SELECT id, user_id, name, created_at, updated_at FROM playlists WHERE id = $1`
	p := &models.Playlist{}
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.UserID, &p.Name, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

// ListByUser returns all playlists belonging to a user (for sidebar display)
func (r *PlaylistRepository) ListByUser(userID int) ([]models.Playlist, error) {
	query := `SELECT id, user_id, name, created_at, updated_at FROM playlists WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []models.Playlist
	for rows.Next() {
		var p models.Playlist
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		playlists = append(playlists, p)
	}
	return playlists, nil
}

func (r *PlaylistRepository) Rename(id int, name string) error {
	_, err := r.db.Exec(`UPDATE playlists SET name = $1, updated_at = NOW() WHERE id = $2`, name, id)
	return err
}

func (r *PlaylistRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM playlists WHERE id = $1`, id)
	return err
}

// --- Playlist Songs ---

// AddSong appends a song at the end of the playlist (position = current max + 1)
func (r *PlaylistRepository) AddSong(playlistID, songID int) error {
	query := `
		INSERT INTO playlist_songs (playlist_id, song_id, position)
		VALUES ($1, $2, COALESCE((SELECT MAX(position) + 1 FROM playlist_songs WHERE playlist_id = $1), 1))
		ON CONFLICT (playlist_id, song_id) DO NOTHING
	`
	_, err := r.db.Exec(query, playlistID, songID)
	return err
}

func (r *PlaylistRepository) RemoveSong(playlistID, songID int) error {
	_, err := r.db.Exec(`DELETE FROM playlist_songs WHERE playlist_id = $1 AND song_id = $2`, playlistID, songID)
	return err
}

// ListSongs returns songs in a playlist ordered by position
func (r *PlaylistRepository) ListSongs(playlistID, userID int) ([]models.Song, error) {
	query := `
		SELECT s.id, s.title, s.artist, s.album, s.genre, s.duration_seconds, s.source_format,
		       s.flac_path, s.mp3_path, s.cover_path, s.transcode_status, s.uploaded_by, s.created_at, s.updated_at,
		       (l.user_id IS NOT NULL) AS is_liked
		FROM songs s
		INNER JOIN playlist_songs ps ON ps.song_id = s.id
		LEFT JOIN liked_songs l ON l.song_id = s.id AND l.user_id = $2
		WHERE ps.playlist_id = $1
		ORDER BY ps.position ASC
	`
	rows, err := r.db.Query(query, playlistID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var s models.Song
				if err := rows.Scan(
			&s.ID, &s.Title, &s.Artist, &s.Album, &s.Genre, &s.DurationSeconds,
			&s.SourceFormat, &s.FlacPath, &s.Mp3Path, &s.CoverPath,
			&s.TranscodeStatus, &s.UploadedBy, &s.CreatedAt, &s.UpdatedAt, &s.IsLiked,
		); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}
	return songs, nil
}

// Reorder updates position for each song_id in the given order (index 0 = position 1, etc.)
func (r *PlaylistRepository) Reorder(playlistID int, songIDsInOrder []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, songID := range songIDsInOrder {
		_, err := tx.Exec(
			`UPDATE playlist_songs SET position = $1 WHERE playlist_id = $2 AND song_id = $3`,
			i+1, playlistID, songID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}