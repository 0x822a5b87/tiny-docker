package conf

import (
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func GenContainerID() (fullID, shortID string) {
	u := uuid.New()
	fullID = strings.ReplaceAll(u.String(), "-", "")
	shortID = fullID[:12]
	return
}

func ExtractNameFromTarPath(tarPath string) string {
	filename := filepath.Base(tarPath)

	suffixes := []string{".tar.gz", ".tgz", ".tar.bz2", ".tar.xz", ".tar"}

	for _, suffix := range suffixes {
		if strings.HasSuffix(strings.ToLower(filename), suffix) {
			filename = strings.TrimSuffix(filename, suffix)
			break
		}
	}

	return filename
}
