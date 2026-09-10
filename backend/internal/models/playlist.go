package models

import "time"

type Playlist struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PlaylistSong struct {
	PlaylistID int       `json:"playlist_id"`
	SongID     int       `json:"song_id"`
	Position   int       `json:"position"`
	AddedAt    time.Time `json:"added_at"`
}