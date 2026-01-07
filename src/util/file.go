package util

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

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

func UnTar(tarPath string, destDir string) error {
	if _, err := exec.Command("tar", "-xvf", tarPath, "-C", destDir).CombinedOutput(); err != nil {
		logrus.Errorf("untar error: {%s}", err.Error())
		return err
	}
	return nil
}
