package tagreader

import (
	"os"

	"github.com/dhowden/tag"
)

type Metadata struct {
	Title      string
	Artist     string
	Album      string
	Genre      string
	CoverData  []byte // raw image bytes, nil if no embedded cover
	CoverExt   string // "jpg" or "png"
}

// Read extracts ID3/FLAC tags from an audio file, including embedded cover
// art if present. Falls back to empty strings/nil for missing fields.
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

	meta := &Metadata{
		Title:  m.Title(),
		Artist: m.Artist(),
		Album:  m.Album(),
		Genre:  m.Genre(),
	}

	if picture := m.Picture(); picture != nil {
		meta.CoverData = picture.Data
		ext := "jpg"
		if picture.MIMEType == "image/png" {
			ext = "png"
		}
		meta.CoverExt = ext
	}

	return meta, nil
}