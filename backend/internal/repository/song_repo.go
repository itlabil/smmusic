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

// FindDuplicate returns an existing song matching the same title+artist
// (case-insensitive), used to detect likely duplicate uploads.
func (r *SongRepository) FindDuplicate(title, artist string) (*models.Song, error) {
	query := `
		SELECT id, title, artist, album, genre, duration_seconds, source_format,
		       flac_path, mp3_path, cover_path, transcode_status, uploaded_by, created_at, updated_at
		FROM songs
		WHERE LOWER(title) = LOWER($1) AND LOWER(artist) = LOWER($2)
		LIMIT 1
	`
	s := &models.Song{}
	err := r.db.QueryRow(query, title, artist).Scan(
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

func (r *SongRepository) ListWithLikedStatus(userID, limit, offset int) ([]models.Song, error) {
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
			&s.TranscodeStatus, &s.UploadedBy, &s.CreatedAt, &s.UpdatedAt, &s.IsLiked, &s.IsInPlaylist,
		); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}
	return songs, nil
}

func (r *SongRepository) Search(userID int, query string, limit int) ([]models.Song, error) {
	sqlQuery := `
		SELECT s.id, s.title, s.artist, s.album, s.genre, s.duration_seconds, s.source_format,
		       s.flac_path, s.mp3_path, s.cover_path, s.transcode_status, s.uploaded_by, s.created_at, s.updated_at,
		       (l.user_id IS NOT NULL) AS is_liked,
		       EXISTS (
		           SELECT 1 FROM playlist_songs ps
		           INNER JOIN playlists p ON p.id = ps.playlist_id
		           WHERE ps.song_id = s.id AND p.user_id = $1
		       ) AS is_in_playlist
		FROM songs s
		LEFT JOIN liked_songs l ON l.song_id = s.id AND l.user_id = $1
		WHERE s.search_vector @@ plainto_tsquery('simple', $2)
		ORDER BY ts_rank(s.search_vector, plainto_tsquery('simple', $2)) DESC
		LIMIT $3
	`
	rows, err := r.db.Query(sqlQuery, userID, query, limit)
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

// UpdateSingleFormatPath sets only flac_path or mp3_path, leaving the other
// column untouched. Used when merging a newly uploaded file into an existing
// duplicate song.
func (r *SongRepository) UpdateSingleFormatPath(id int, format, path string) error {
	var query string
	if format == "flac" {
		query = `UPDATE songs SET flac_path = $1, updated_at = NOW() WHERE id = $2`
	} else {
		query = `UPDATE songs SET mp3_path = $1, updated_at = NOW() WHERE id = $2`
	}
	_, err := r.db.Exec(query, path, id)
	return err
}

func (r *SongRepository) UpdateCoverPath(id int, coverPath string) error {
	query := `UPDATE songs SET cover_path = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, coverPath, id)
	return err
}

func (r *SongRepository) UpdateMetadata(id int, title, artist string, album, genre *string) error {
	query := `
		UPDATE songs
		SET title = $1, artist = $2, album = $3, genre = $4, updated_at = NOW()
		WHERE id = $5
	`
	_, err := r.db.Exec(query, title, artist, album, genre, id)
	return err
}

func (r *SongRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM songs WHERE id = $1`, id)
	return err
}