package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itlabil/smmusic/backend/internal/service"
)

type InteractionHandler struct {
	service *service.InteractionService
}

func NewInteractionHandler(s *service.InteractionService) *InteractionHandler {
	return &InteractionHandler{service: s}
}

func (h *InteractionHandler) Like(c *gin.Context) {
	userID := c.GetInt("user_id")
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid song id"})
		return
	}

	if err := h.service.LikeSong(userID, songID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "song liked"})
}

func (h *InteractionHandler) Unlike(c *gin.Context) {
	userID := c.GetInt("user_id")
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid song id"})
		return
	}

	if err := h.service.UnlikeSong(userID, songID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "song unliked"})
}

func (h *InteractionHandler) ListLiked(c *gin.Context) {
	userID := c.GetInt("user_id")
	songs, err := h.service.ListLikedSongs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"songs": songs})
}

func (h *InteractionHandler) RecordPlay(c *gin.Context) {
	userID := c.GetInt("user_id")
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid song id"})
		return
	}

	if err := h.service.RecordPlay(userID, songID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "play recorded"})
}

func (h *InteractionHandler) ListRecentlyPlayed(c *gin.Context) {
	userID := c.GetInt("user_id")
	songs, err := h.service.ListRecentlyPlayed(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"songs": songs})
}

func (h *InteractionHandler) LikedSummary(c *gin.Context) {
	userID := c.GetInt("user_id")
	summary, err := h.service.GetLikedSongsSummary(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"summary": summary})
}