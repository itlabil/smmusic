package models

import "time"

type LikedSong struct {
	UserID   int       `json:"user_id"`
	SongID   int       `json:"song_id"`
	LikedAt  time.Time `json:"liked_at"`
}

type PlayHistory struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	SongID   int       `json:"song_id"`
	PlayedAt time.Time `json:"played_at"`
}