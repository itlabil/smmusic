package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itlabil/smmusic/backend/internal/service"
)

type SongHandler struct {
	songService *service.SongService
}

func NewSongHandler(songService *service.SongService) *SongHandler {
	return &SongHandler{songService: songService}
}

func (h *SongHandler) Upload(c *gin.Context) {
	userID := c.GetInt("user_id")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required (multipart field name: 'file')"})
		return
	}

	// Limit upload size to 200MB (generous for FLAC files)
	const maxUploadSize = 200 << 20
	if fileHeader.Size > maxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 200MB)"})
		return
	}

	song, err := h.songService.Upload(fileHeader, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"song": song})
}

func (h *SongHandler) List(c *gin.Context) {
	songs, err := h.songService.List(50, 0) // basic pagination default, refined later
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"songs": songs})
}

func (h *SongHandler) Stream(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid song id"})
		return
	}

	quality := c.DefaultQuery("quality", "standard") // "standard" | "high"

	filePath, contentType, err := h.songService.GetStreamPath(songID, quality)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", contentType)
	c.Header("Accept-Ranges", "bytes")
	// http.ServeFile handles Range requests automatically (seek/skip support)
	c.File(filePath)
	_ = contentType
}