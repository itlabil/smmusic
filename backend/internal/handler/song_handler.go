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
	userID := c.GetInt("user_id")
	songs, err := h.songService.List(userID, 50, 0) // basic pagination default, refined later
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

func (h *SongHandler) Cover(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid song id"})
		return
	}

	coverPath, err := h.songService.GetCoverPath(songID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.File(coverPath)
}

type updateSongRequest struct {
	Title  string  `json:"title" binding:"required"`
	Artist string  `json:"artist" binding:"required"`
	Album  *string `json:"album"`
	Genre  *string `json:"genre"`
}

func (h *SongHandler) Update(c *gin.Context) {
	userID := c.GetInt("user_id")
	role := c.GetString("role")
	isAdmin := role == "admin"

	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid song id"})
		return
	}

	var req updateSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and artist are required"})
		return
	}

	if err := h.songService.UpdateMetadata(songID, userID, isAdmin, req.Title, req.Artist, req.Album, req.Genre); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "song updated"})
}

func (h *SongHandler) UploadCover(c *gin.Context) {
	userID := c.GetInt("user_id")
	role := c.GetString("role")
	isAdmin := role == "admin"

	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid song id"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required (multipart field name: 'file')"})
		return
	}

	if err := h.songService.UploadCover(songID, userID, isAdmin, fileHeader); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cover uploaded"})
}

func (h *SongHandler) Search(c *gin.Context) {
	userID := c.GetInt("user_id")
	query := c.Query("q")

	if query == "" {
		c.JSON(http.StatusOK, gin.H{"songs": []interface{}{}})
		return
	}

	songs, err := h.songService.Search(userID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"songs": songs})
}