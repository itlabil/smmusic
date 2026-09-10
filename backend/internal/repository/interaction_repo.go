package repository

import (
	"database/sql"

	"github.com/itlabil/smmusic/backend/internal/models"
)

type InteractionRepository struct {
	db *sql.DB
}

func NewInteractionRepository(db *sql.DB) *InteractionRepository {
	return &InteractionRepository{db: db}
}

// --- Liked Songs ---

func (r *InteractionRepository) LikeSong(userID, songID int) error {
	query := `
		INSERT INTO liked_songs (user_id, song_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, song_id) DO NOTHING
	`
	_, err := r.db.Exec(query, userID, songID)
	return err
}

func (r *InteractionRepository) UnlikeSong(userID, songID int) error {
	_, err := r.db.Exec(`DELETE FROM liked_songs WHERE user_id = $1 AND song_id = $2`, userID, songID)
	return err
}

func (r *InteractionRepository) ListLikedSongs(userID int) ([]models.Song, error) {
	query := `
		SELECT s.id, s.title, s.artist, s.album, s.genre, s.duration_seconds, s.source_format,
		       s.flac_path, s.mp3_path, s.cover_path, s.transcode_status, s.uploaded_by, s.created_at, s.updated_at,
		       EXISTS (
		           SELECT 1 FROM playlist_songs ps
		           INNER JOIN playlists p ON p.id = ps.playlist_id
		           WHERE ps.song_id = s.id AND p.user_id = $1
		       ) AS is_in_playlist
		FROM songs s
		INNER JOIN liked_songs l ON l.song_id = s.id
		WHERE l.user_id = $1
		ORDER BY l.liked_at DESC
	`
	rows, err := r.db.Query(query, userID)
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
			&s.TranscodeStatus, &s.UploadedBy, &s.CreatedAt, &s.UpdatedAt, &s.IsInPlaylist,
		); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}
	return songs, nil
}

// --- Play History ---
func (r *InteractionRepository) RecordPlay(userID, songID int) error {
	_, err := r.db.Exec(`INSERT INTO play_history (user_id, song_id) VALUES ($1, $2)`, userID, songID)
	return err
}

func (r *InteractionRepository) ListRecentlyPlayed(userID, limit int) ([]models.Song, error) {
	query := `
		SELECT s.id, s.title, s.artist, s.album, s.genre, s.duration_seconds, s.source_format,
		       s.flac_path, s.mp3_path, s.cover_path, s.transcode_status, s.uploaded_by, s.created_at, s.updated_at,
		       (l.user_id IS NOT NULL) AS is_liked,
		       EXISTS (
		           SELECT 1 FROM playlist_songs ps
		           INNER JOIN playlists p ON p.id = ps.playlist_id
		           WHERE ps.song_id = s.id AND p.user_id = $1
		       ) AS is_in_playlist
		FROM songs s
		INNER JOIN (
			SELECT DISTINCT ON (song_id) song_id, played_at
			FROM play_history
			WHERE user_id = $1
			ORDER BY song_id, played_at DESC
		) latest ON latest.song_id = s.id
		LEFT JOIN liked_songs l ON l.song_id = s.id AND l.user_id = $1
		ORDER BY latest.played_at DESC
		LIMIT $2
	`
	rows, err := r.db.Query(query, userID, limit)
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
			&s.TranscodeStatus, &s.UploadedBy, &s.CreatedAt, &s.UpdatedAt, &s.IsLiked, &s.IsInPlaylist,
		); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}
	return songs, nil
}