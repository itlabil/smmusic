package slug

import (
	"regexp"
	"strings"
)

var nonAlnumRegex = regexp.MustCompile(`[^a-z0-9]+`)

// Generate turns "AC/DC" or "Sheila On 7" into "ac-dc" / "sheila-on-7",
// used as a safe folder name for storage/audio/{artist}/
func Generate(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonAlnumRegex.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "unknown-artist"
	}
	return s
}