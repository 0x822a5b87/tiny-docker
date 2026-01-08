package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/sirupsen/logrus"
)

func NullFile() (*os.File, error) {
	return os.OpenFile(constant.NullFilePath, os.O_RDWR, 0666)
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

func Tar(dstFile string, srcPath string) error {
	data, err := exec.Command("tar", "-cvf", dstFile, "-C", srcPath, ".").CombinedOutput()
	logrus.Debug(string(data))
	if err != nil {
		logrus.Errorf("tar error: {%s}", err.Error())
		logrus.Errorf("tar error: dst = {%s}, src = {%s}, error = {%s}", dstFile, srcPath, err.Error())
		logrus.Errorf("tar error info: {%s}", string(data))
		return err
	}
	return nil
}

func UnTar(srcFile string, dstPath string) error {
	data, err := exec.Command("tar", "-xvf", srcFile, "-C", dstPath).CombinedOutput()
	logrus.Debug(string(data))
	if err != nil {
		logrus.Errorf("untar error: {%s}", err.Error())
		logrus.Errorf("untar error: src = {%s}, dst = {%s}, error = {%s}", srcFile, dstPath, err.Error())
		logrus.Errorf("untar error info: {%s}", string(data))
		return err
	}
	return nil
}
