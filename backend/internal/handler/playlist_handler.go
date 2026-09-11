package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itlabil/smmusic/backend/internal/service"
)

type PlaylistHandler struct {
	service *service.PlaylistService
}

func NewPlaylistHandler(s *service.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{service: s}
}

type createPlaylistRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *PlaylistHandler) Create(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req createPlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	playlist, err := h.service.Create(userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"playlist": playlist})
}

func (h *PlaylistHandler) List(c *gin.Context) {
	userID := c.GetInt("user_id")
	playlists, err := h.service.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"playlists": playlists})
}

type renamePlaylistRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *PlaylistHandler) Rename(c *gin.Context) {
	userID := c.GetInt("user_id")
	playlistID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist id"})
		return
	}

	var req renamePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	if err := h.service.Rename(playlistID, userID, req.Name); err != nil {
		handlePlaylistError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "playlist renamed"})
}

func (h *PlaylistHandler) Delete(c *gin.Context) {
	userID := c.GetInt("user_id")
	playlistID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist id"})
		return
	}

	if err := h.service.Delete(playlistID, userID); err != nil {
		handlePlaylistError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "playlist deleted"})
}

func (h *PlaylistHandler) ListSongs(c *gin.Context) {
	userID := c.GetInt("user_id")
	playlistID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist id"})
		return
	}

	songs, err := h.service.ListSongs(playlistID, userID)
	if err != nil {
		handlePlaylistError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"songs": songs})
}

type addSongRequest struct {
	SongID int `json:"song_id" binding:"required"`
}

func (h *PlaylistHandler) AddSong(c *gin.Context) {
	userID := c.GetInt("user_id")
	playlistID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist id"})
		return
	}

	var req addSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "song_id is required"})
		return
	}

	if err := h.service.AddSong(playlistID, userID, req.SongID); err != nil {
		handlePlaylistError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "song added to playlist"})
}

func (h *PlaylistHandler) RemoveSong(c *gin.Context) {
	userID := c.GetInt("user_id")
	playlistID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist id"})
		return
	}
	songID, err := strconv.Atoi(c.Param("songId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid song id"})
		return
	}

	if err := h.service.RemoveSong(playlistID, userID, songID); err != nil {
		handlePlaylistError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "song removed from playlist"})
}

type reorderRequest struct {
	SongIDs []int `json:"song_ids" binding:"required"`
}

func (h *PlaylistHandler) Reorder(c *gin.Context) {
	userID := c.GetInt("user_id")
	playlistID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist id"})
		return
	}

	var req reorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "song_ids array is required"})
		return
	}

	if err := h.service.Reorder(playlistID, userID, req.SongIDs); err != nil {
		handlePlaylistError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "playlist reordered"})
}

func handlePlaylistError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrPlaylistNotOwned) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this playlist"})
		return
	}
	if err.Error() == "playlist not found" {
		c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

func (h *PlaylistHandler) GetDetail(c *gin.Context) {
	userID := c.GetInt("user_id")
	playlistID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist id"})
		return
	}

	detail, err := h.service.GetDetail(playlistID, userID)
	if err != nil {
		handlePlaylistError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"playlist": detail})
}