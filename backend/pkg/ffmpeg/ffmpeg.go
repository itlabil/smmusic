package ffmpeg

import (
	"fmt"
	"os/exec"
)

// TranscodeToMP3 converts a FLAC file to MP3 at the given bitrate (e.g. "256k").
// Requires ffmpeg to be installed and available in PATH.
func TranscodeToMP3(inputPath, outputPath, bitrate string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputPath,
		"-codec:a", "libmp3lame",
		"-b:a", bitrate,
		"-y", // overwrite output if exists
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg transcode failed: %w (output: %s)", err, string(output))
	}

	return nil
}