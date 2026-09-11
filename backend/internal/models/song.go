package models

import "time"

type Song struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Artist           string    `json:"artist"`
	Album            *string   `json:"album,omitempty"`
	Genre            *string   `json:"genre,omitempty"`
	DurationSeconds  *int      `json:"duration_seconds,omitempty"`
	SourceFormat     string    `json:"source_format"`
	FlacPath         *string   `json:"flac_path,omitempty"`
	Mp3Path          *string   `json:"mp3_path,omitempty"`
	CoverPath        *string   `json:"cover_path,omitempty"`
	TranscodeStatus  string    `json:"transcode_status"`
	UploadedBy       *int      `json:"uploaded_by,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	IsLiked          bool      `json:"is_liked"`
	IsInPlaylist     bool      `json:"is_in_playlist"`
	AddedAt          *time.Time `json:"added_at,omitempty"`
}