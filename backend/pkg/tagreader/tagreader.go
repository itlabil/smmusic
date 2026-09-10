package tagreader

import (
	"os"

	"github.com/dhowden/tag"
)

type Metadata struct {
	Title  string
	Artist string
	Album  string
	Genre  string
}

// Read extracts ID3/FLAC tags from an audio file. Falls back to empty
// strings for any field not present in the file's metadata.
func Read(filePath string) (*Metadata, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	return &Metadata{
		Title:  m.Title(),
		Artist: m.Artist(),
		Album:  m.Album(),
		Genre:  m.Genre(),
	}, nil
}