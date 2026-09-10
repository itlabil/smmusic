package repository

import (
	"database/sql"
	"errors"

	"github.com/itlabil/smmusic/backend/internal/models"
)

type SongRepository struct {
	db *sql.DB
}

func NewSongRepository(db *sql.DB) *SongRepository {
	return &SongRepository{db: db}
}

func (r *SongRepository) Create(s *models.Song) error {
	query := `
		INSERT INTO songs (title, artist, album, genre, duration_seconds, source_format, flac_path, mp3_path, cover_path, transcode_status, uploaded_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		s.Title, s.Artist, s.Album, s.Genre, s.DurationSeconds,
		s.SourceFormat, s.FlacPath, s.Mp3Path, s.CoverPath,
		s.TranscodeStatus, s.UploadedBy,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *SongRepository) FindByID(id int) (*models.Song, error) {
	query := `
		SELECT id, title, artist, album, genre, duration_seconds, source_format,
		       flac_path, mp3_path, cover_path, transcode_status, uploaded_by, created_at, updated_at
		FROM songs WHERE id = $1
	`
	s := &models.Song{}
	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.Title, &s.Artist, &s.Album, &s.Genre, &s.DurationSeconds,
		&s.SourceFormat, &s.FlacPath, &s.Mp3Path, &s.CoverPath,
		&s.TranscodeStatus, &s.UploadedBy, &s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

// List returns songs ordered by newest first (basic library view)
func (r *SongRepository) List(limit, offset int) ([]models.Song, error) {
	query := `
		SELECT id, title, artist, album, genre, duration_seconds, source_format,
		       flac_path, mp3_path, cover_path, transcode_status, uploaded_by, created_at, updated_at
		FROM songs ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
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
			&s.TranscodeStatus, &s.UploadedBy, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}
	return songs, nil
}

// Search using Postgres full-text search on search_vector
func (r *SongRepository) Search(query string, limit int) ([]models.Song, error) {
	sqlQuery := `
		SELECT id, title, artist, album, genre, duration_seconds, source_format,
		       flac_path, mp3_path, cover_path, transcode_status, uploaded_by, created_at, updated_at
		FROM songs
		WHERE search_vector @@ plainto_tsquery('simple', $1)
		ORDER BY ts_rank(search_vector, plainto_tsquery('simple', $1)) DESC
		LIMIT $2
	`
	rows, err := r.db.Query(sqlQuery, query, limit)
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
			&s.TranscodeStatus, &s.UploadedBy, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}
	return songs, nil
}

// UpdateTranscodeStatus is called by the background transcode job
func (r *SongRepository) UpdateTranscodeStatus(id int, status string, mp3Path *string) error {
	query := `UPDATE songs SET transcode_status = $1, mp3_path = COALESCE($2, mp3_path), updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(query, status, mp3Path, id)
	return err
}

func (r *SongRepository) UpdatePaths(id int, flacPath, mp3Path *string) error {
	query := `UPDATE songs SET flac_path = $1, mp3_path = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(query, flacPath, mp3Path, id)
	return err
}

func (r *SongRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM songs WHERE id = $1`, id)
	return err
}

// ListWithLikedStatus returns songs with an is_liked flag for the given user
func (r *SongRepository) ListWithLikedStatus(userID, limit, offset int) ([]models.Song, error) {
	query := `
		SELECT s.id, s.title, s.artist, s.album, s.genre, s.duration_seconds, s.source_format,
		       s.flac_path, s.mp3_path, s.cover_path, s.transcode_status, s.uploaded_by, s.created_at, s.updated_at,
		       (l.user_id IS NOT NULL) AS is_liked
		FROM songs s
		LEFT JOIN liked_songs l ON l.song_id = s.id AND l.user_id = $1
		ORDER BY s.created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, userID, limit, offset)
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