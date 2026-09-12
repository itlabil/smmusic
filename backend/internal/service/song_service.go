package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/itlabil/smmusic/backend/internal/config"
	"github.com/itlabil/smmusic/backend/internal/models"
	"github.com/itlabil/smmusic/backend/internal/repository"
	"github.com/itlabil/smmusic/backend/pkg/ffmpeg"
	"github.com/itlabil/smmusic/backend/pkg/slug"
	"github.com/itlabil/smmusic/backend/pkg/tagreader"
)

// pendingUpload holds everything needed to finish processing an upload after
// the user resolves a duplicate-song prompt (merge into existing vs. new song).
type pendingUpload struct {
	TempPath        string
	Ext             string
	SourceFormat    string
	Title           string
	Artist          string
	Album           *string
	Genre           *string
	CoverData       []byte
	CoverExt        string
	DurationSeconds int
	UploadedBy      int
	ExistingSongID  *int
}

type DuplicateInfo struct {
	ExistingSong *models.Song `json:"existing_song"`
	CanMerge     bool         `json:"can_merge"`
	UploadToken  string       `json:"upload_token"`
}

type UploadResult struct {
	Song      *models.Song   `json:"song,omitempty"`
	Duplicate *DuplicateInfo `json:"duplicate,omitempty"`
}

type SongService struct {
	songRepo         *repository.SongRepository
	cfg              *config.Config
	transcodeService *TranscodeService

	mu             sync.Mutex
	pendingUploads map[string]*pendingUpload
}

func NewSongService(songRepo *repository.SongRepository, cfg *config.Config, transcodeService *TranscodeService) *SongService {
	return &SongService{
		songRepo:         songRepo,
		cfg:              cfg,
		transcodeService: transcodeService,
		pendingUploads:   make(map[string]*pendingUpload),
	}
}

func generateUploadToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Upload processes a newly uploaded audio file. If a song with the same
// title+artist already exists, it returns a Duplicate result instead of
// creating a new song, so the caller can decide (via ConfirmUpload) whether
// to merge the file into the existing song or upload it as a new song.
func (s *SongService) Upload(fileHeader *multipart.FileHeader, uploadedBy int) (*UploadResult, error) {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".flac" && ext != ".mp3" {
		return nil, fmt.Errorf("unsupported file type: %s (only .flac and .mp3 allowed)", ext)
	}

	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("upload-%d-%d%s", uploadedBy, os.Getpid(), ext))
	if err := saveMultipartFile(fileHeader, tempPath); err != nil {
		return nil, err
	}

	meta, err := tagreader.Read(tempPath)
	if err != nil {
		os.Remove(tempPath)
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

	durationSeconds, err := ffmpeg.GetDuration(tempPath)
	if err != nil {
		durationSeconds = 0
	}

	var albumPtr, genrePtr *string
	if meta.Album != "" {
		albumPtr = &meta.Album
	}
	if meta.Genre != "" {
		genrePtr = &meta.Genre
	}

	pending := &pendingUpload{
		TempPath:        tempPath,
		Ext:             ext,
		SourceFormat:    sourceFormat,
		Title:           meta.Title,
		Artist:          meta.Artist,
		Album:           albumPtr,
		Genre:           genrePtr,
		CoverData:       meta.CoverData,
		CoverExt:        meta.CoverExt,
		DurationSeconds: durationSeconds,
		UploadedBy:      uploadedBy,
	}

	existing, err := s.songRepo.FindDuplicate(meta.Title, meta.Artist)
	if err != nil {
		os.Remove(tempPath)
		return nil, err
	}

	if existing != nil {
		canMerge := (sourceFormat == "flac" && existing.FlacPath == nil) ||
			(sourceFormat == "mp3" && existing.Mp3Path == nil)

		pending.ExistingSongID = &existing.ID

		token, err := generateUploadToken()
		if err != nil {
			os.Remove(tempPath)
			return nil, err
		}

		s.mu.Lock()
		s.pendingUploads[token] = pending
		s.mu.Unlock()

		return &UploadResult{
			Duplicate: &DuplicateInfo{
				ExistingSong: existing,
				CanMerge:     canMerge,
				UploadToken:  token,
			},
		}, nil
	}

	song, err := s.finalizeNewSong(pending)
	if err != nil {
		return nil, err
	}
	return &UploadResult{Song: song}, nil
}

// ConfirmUpload is called after the user resolves a duplicate-song prompt.
// action "new" creates a separate song as usual; action "merge" attaches the
// uploaded file to the existing matching song instead.
func (s *SongService) ConfirmUpload(token, action string) (*models.Song, error) {
	s.mu.Lock()
	pending, ok := s.pendingUploads[token]
	if ok {
		delete(s.pendingUploads, token)
	}
	s.mu.Unlock()

	if !ok {
		return nil, fmt.Errorf("upload token not found or expired")
	}

	if action == "merge" {
		if pending.ExistingSongID == nil {
			return nil, fmt.Errorf("no existing song to merge into")
		}
		return s.mergeIntoExisting(pending, *pending.ExistingSongID)
	}

	return s.finalizeNewSong(pending)
}

func (s *SongService) finalizeNewSong(pending *pendingUpload) (*models.Song, error) {
	defer os.Remove(pending.TempPath)

	song := &models.Song{
		Title:           pending.Title,
		Artist:          pending.Artist,
		Album:           pending.Album,
		Genre:           pending.Genre,
		DurationSeconds: &pending.DurationSeconds,
		SourceFormat:    pending.SourceFormat,
		TranscodeStatus: "done",
		UploadedBy:      &pending.UploadedBy,
	}
	if pending.SourceFormat == "flac" {
		song.TranscodeStatus = "pending"
	}

	if err := s.songRepo.Create(song); err != nil {
		return nil, err
	}

	artistFolder := slug.Generate(song.Artist)
	destDir := filepath.Join(s.cfg.StorageAudioPath, artistFolder)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, err
	}

	destPath := filepath.Join(destDir, fmt.Sprintf("%d%s", song.ID, pending.Ext))
	if err := copyFile(pending.TempPath, destPath); err != nil {
		return nil, err
	}

	if pending.SourceFormat == "flac" {
		song.FlacPath = &destPath
	} else {
		song.Mp3Path = &destPath
	}
	if err := s.songRepo.UpdatePaths(song.ID, song.FlacPath, song.Mp3Path); err != nil {
		return nil, err
	}

	if len(pending.CoverData) > 0 {
		coverDir := filepath.Join(s.cfg.StorageCoverPath, artistFolder)
		if err := os.MkdirAll(coverDir, 0755); err != nil {
			return nil, err
		}
		coverPath := filepath.Join(coverDir, fmt.Sprintf("%d.%s", song.ID, pending.CoverExt))
		if err := os.WriteFile(coverPath, pending.CoverData, 0644); err != nil {
			return nil, err
		}
		song.CoverPath = &coverPath
		if err := s.songRepo.UpdateCoverPath(song.ID, coverPath); err != nil {
			return nil, err
		}
	}

	if pending.SourceFormat == "flac" {
		s.transcodeService.Enqueue(song.ID, destPath)
	}

	return song, nil
}

// mergeIntoExisting attaches the uploaded file's format to an already-existing
// song instead of creating a duplicate row.
func (s *SongService) mergeIntoExisting(pending *pendingUpload, existingID int) (*models.Song, error) {
	defer os.Remove(pending.TempPath)

	existing, err := s.songRepo.FindByID(existingID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("existing song not found")
	}

	artistFolder := slug.Generate(existing.Artist)
	destDir := filepath.Join(s.cfg.StorageAudioPath, artistFolder)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, err
	}

	destPath := filepath.Join(destDir, fmt.Sprintf("%d%s", existing.ID, pending.Ext))
	if err := copyFile(pending.TempPath, destPath); err != nil {
		return nil, err
	}

	if err := s.songRepo.UpdateSingleFormatPath(existing.ID, pending.SourceFormat, destPath); err != nil {
		return nil, err
	}

	if existing.CoverPath == nil && len(pending.CoverData) > 0 {
		coverDir := filepath.Join(s.cfg.StorageCoverPath, artistFolder)
		if err := os.MkdirAll(coverDir, 0755); err == nil {
			coverPath := filepath.Join(coverDir, fmt.Sprintf("%d.%s", existing.ID, pending.CoverExt))
			if err := os.WriteFile(coverPath, pending.CoverData, 0644); err == nil {
				_ = s.songRepo.UpdateCoverPath(existing.ID, coverPath)
			}
		}
	}

	return s.songRepo.FindByID(existing.ID)
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
	songs, err := s.songRepo.ListWithLikedStatus(userID, limit, offset)
	if err != nil {
		return nil, err
	}
	for i := range songs {
		songs[i].AddedAt = &songs[i].CreatedAt
	}
	return songs, nil
}

func (s *SongService) Search(userID int, query string) ([]models.Song, error) {
	return s.songRepo.Search(userID, query, 30)
}

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

func (s *SongService) GetCoverPath(songID int) (string, error) {
	song, err := s.songRepo.FindByID(songID)
	if err != nil {
		return "", err
	}
	if song == nil || song.CoverPath == nil {
		return "", fmt.Errorf("no cover art available")
	}
	return *song.CoverPath, nil
}

func (s *SongService) UpdateMetadata(songID, requestingUserID int, isAdmin bool, title, artist string, album, genre *string) error {
	song, err := s.songRepo.FindByID(songID)
	if err != nil {
		return err
	}
	if song == nil {
		return fmt.Errorf("song not found")
	}

	if !isAdmin && (song.UploadedBy == nil || *song.UploadedBy != requestingUserID) {
		return fmt.Errorf("you do not have permission to edit this song")
	}

	return s.songRepo.UpdateMetadata(songID, title, artist, album, genre)
}

func (s *SongService) UploadCover(songID, requestingUserID int, isAdmin bool, fileHeader *multipart.FileHeader) error {
	song, err := s.songRepo.FindByID(songID)
	if err != nil {
		return err
	}
	if song == nil {
		return fmt.Errorf("song not found")
	}

	if !isAdmin && (song.UploadedBy == nil || *song.UploadedBy != requestingUserID) {
		return fmt.Errorf("you do not have permission to edit this song")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return fmt.Errorf("unsupported image type: %s (only .jpg and .png allowed)", ext)
	}
	if ext == ".jpeg" {
		ext = ".jpg"
	}

	artistFolder := slug.Generate(song.Artist)
	coverDir := filepath.Join(s.cfg.StorageCoverPath, artistFolder)
	if err := os.MkdirAll(coverDir, 0755); err != nil {
		return err
	}

	coverPath := filepath.Join(coverDir, fmt.Sprintf("%d%s", songID, ext))
	if err := saveMultipartFile(fileHeader, coverPath); err != nil {
		return err
	}

	return s.songRepo.UpdateCoverPath(songID, coverPath)
}

func (s *SongService) Delete(songID, requestingUserID int, isAdmin bool) error {
	song, err := s.songRepo.FindByID(songID)
	if err != nil {
		return err
	}
	if song == nil {
		return fmt.Errorf("song not found")
	}

	if !isAdmin && (song.UploadedBy == nil || *song.UploadedBy != requestingUserID) {
		return fmt.Errorf("you do not have permission to delete this song")
	}

	if err := s.songRepo.Delete(songID); err != nil {
		return err
	}

	if song.FlacPath != nil {
		_ = os.Remove(*song.FlacPath)
	}
	if song.Mp3Path != nil {
		_ = os.Remove(*song.Mp3Path)
	}
	if song.CoverPath != nil {
		_ = os.Remove(*song.CoverPath)
	}

	return nil
}