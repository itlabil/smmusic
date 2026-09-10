package service

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/itlabil/smmusic/backend/internal/repository"
	"github.com/itlabil/smmusic/backend/pkg/ffmpeg"
)

type transcodeJob struct {
	SongID   int
	FlacPath string
}

type TranscodeService struct {
	songRepo    *repository.SongRepository
	jobs        chan transcodeJob
	maxWorkers  int
}

// NewTranscodeService starts a fixed pool of worker goroutines (maxWorkers)
// that process transcode jobs from a buffered channel, so CPU usage stays
// bounded even if many FLAC files are uploaded at once.
func NewTranscodeService(songRepo *repository.SongRepository, maxWorkers int) *TranscodeService {
	ts := &TranscodeService{
		songRepo:   songRepo,
		jobs:       make(chan transcodeJob, 100), // buffer up to 100 pending jobs
		maxWorkers: maxWorkers,
	}
	ts.startWorkers()
	return ts
}

func (ts *TranscodeService) startWorkers() {
	for i := 0; i < ts.maxWorkers; i++ {
		go ts.worker(i)
	}
}

func (ts *TranscodeService) worker(id int) {
	for job := range ts.jobs {
		log.Printf("[transcode worker %d] processing song_id=%d", id, job.SongID)

		if err := ts.songRepo.UpdateTranscodeStatus(job.SongID, "processing", nil); err != nil {
			log.Printf("[transcode worker %d] failed to update status to processing: %v", id, err)
			continue
		}

		ext := filepath.Ext(job.FlacPath)
		mp3Path := strings.TrimSuffix(job.FlacPath, ext) + ".mp3"

		if err := ffmpeg.TranscodeToMP3(job.FlacPath, mp3Path, "256k"); err != nil {
			log.Printf("[transcode worker %d] transcode failed for song_id=%d: %v", id, job.SongID, err)
			_ = ts.songRepo.UpdateTranscodeStatus(job.SongID, "failed", nil)
			continue
		}

		if err := ts.songRepo.UpdateTranscodeStatus(job.SongID, "done", &mp3Path); err != nil {
			log.Printf("[transcode worker %d] failed to update status to done: %v", id, err)
			continue
		}

		log.Printf("[transcode worker %d] finished song_id=%d -> %s", id, job.SongID, mp3Path)
	}
}

// Enqueue adds a song to the transcode queue. Non-blocking unless the buffer is full.
func (ts *TranscodeService) Enqueue(songID int, flacPath string) {
	ts.jobs <- transcodeJob{SongID: songID, FlacPath: flacPath}
	log.Printf("enqueued transcode job for song_id=%d", songID)
}

var _ = fmt.Sprintf // placeholder to avoid unused import if trimmed later