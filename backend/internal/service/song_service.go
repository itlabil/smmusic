package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/itlabil/smmusic/backend/internal/config"
	"github.com/itlabil/smmusic/backend/internal/models"
	"github.com/itlabil/smmusic/backend/internal/repository"
	"github.com/itlabil/smmusic/backend/pkg/slug"
	"github.com/itlabil/smmusic/backend/pkg/tagreader"
)

type SongService struct {
	songRepo         *repository.SongRepository
	cfg              *config.Config
	transcodeService *TranscodeService
}

func NewSongService(songRepo *repository.SongRepository, cfg *config.Config, transcodeService *TranscodeService) *SongService {
	return &SongService{songRepo: songRepo, cfg: cfg, transcodeService: transcodeService}
}

// Upload saves the uploaded file to storage/audio/{artist_slug}/{song_id}.{ext},
// extracts metadata from tags, and inserts the songs row. If the source is FLAC,
// transcode_status is left "pending" for the background job (Step 9) to pick up.
func (s *SongService) Upload(fileHeader *multipart.FileHeader, uploadedBy int) (*models.Song, error) {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename)) // ".flac" or ".mp3"
	if ext != ".flac" && ext != ".mp3" {
		return nil, fmt.Errorf("unsupported file type: %s (only .flac and .mp3 allowed)", ext)
	}

	// Save to a temp path first so tagreader can read it
	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("upload-%d%s", uploadedBy, ext))
	if err := saveMultipartFile(fileHeader, tempPath); err != nil {
		return nil, err
	}
	defer os.Remove(tempPath)

	meta, err := tagreader.Read(tempPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio metadata: %w", err)
	}
	if meta.Title == "" {
		meta.Title = strings.TrimSuffix(fileHeader.Filename, ext)
	}
	if meta.Artist == "" {
		meta.Artist = "Unknown Artist"
	}

	sourceFormat := "mp3"
	if ext == ".flac" {
		sourceFormat = "flac"
	}

	song := &models.Song{
		Title:           meta.Title,
		Artist:          meta.Artist,
		SourceFormat:    sourceFormat,
		TranscodeStatus: "done", // default; overridden below if flac
		UploadedBy:      &uploadedBy,
	}
	if meta.Album != "" {
		song.Album = &meta.Album
	}
	if meta.Genre != "" {
		song.Genre = &meta.Genre
	}

	if sourceFormat == "flac" {
		song.TranscodeStatus = "pending"
	}

	// Insert first to get the song ID (used as filename)
	if err := s.songRepo.Create(song); err != nil {
		return nil, err
	}

	artistFolder := slug.Generate(song.Artist)
	destDir := filepath.Join(s.cfg.StorageAudioPath, artistFolder)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, err
	}

	destPath := filepath.Join(destDir, fmt.Sprintf("%d%s", song.ID, ext))
	if err := copyFile(tempPath, destPath); err != nil {
		return nil, err
	}

		if sourceFormat == "flac" {
		song.FlacPath = &destPath
	} else {
		song.Mp3Path = &destPath
	}

	if err := s.songRepo.UpdatePaths(song.ID, song.FlacPath, song.Mp3Path); err != nil {
		return nil, err
	}

	if sourceFormat == "flac" {
		s.transcodeService.Enqueue(song.ID, destPath)
	}

	return song, nil
}

func saveMultipartFile(fh *multipart.FileHeader, dest string) error {
	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func (s *SongService) List(userID, limit, offset int) ([]models.Song, error) {
	return s.songRepo.ListWithLikedStatus(userID, limit, offset)
}

// GetStreamPath resolves which file to serve based on requested quality
// and actual availability. Falls back to MP3 if FLAC isn't available.
func (s *SongService) GetStreamPath(songID int, quality string) (filePath string, contentType string, err error) {
	song, err := s.songRepo.FindByID(songID)
	if err != nil {
		return "", "", err
	}
	if song == nil {
		return "", "", fmt.Errorf("song not found")
	}

	if quality == "high" && song.FlacPath != nil {
		return *song.FlacPath, "audio/flac", nil
	}

	if song.Mp3Path == nil {
		return "", "", fmt.Errorf("no playable file available for this song")
	}

	return *song.Mp3Path, "audio/mpeg", nil
}