package util

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/sirupsen/logrus"
)

func ReadAllFilesInDir(dirPath string) (map[string][]byte, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		logrus.Errorf("directory does not exist: %s", dirPath)
		return nil, err
	} else if err != nil {
		logrus.Errorf("failed to stat directory: %s, err: %v", dirPath, err)
		return nil, err
	}

	fileContents := make(map[string][]byte)

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logrus.Warnf("skip file %s, walk error: %v", path, err)
			return nil
		}

		if d.Type().IsRegular() {
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				logrus.Errorf("failed to read file %s, err: %v", path, readErr)
				return nil
			}
			filename := filepath.Base(path)
			fileContents[filename] = content
			logrus.Debugf("successfully read file: %s", path)
		}

		return nil
	})

	if err != nil {
		logrus.Errorf("failed to walk directory %s, err: %v", dirPath, err)
		return nil, err
	}

	return fileContents, nil
}

func GetFdRealPath(f *os.File) (string, error) {
	if f == nil {
		return "", fmt.Errorf("file is nil")
	}
	fd := f.Fd()
	linkPath := filepath.Join("/proc/self/fd", fmt.Sprintf("%d", fd))
	realPath, err := os.Readlink(linkPath)
	if err != nil {
		return "", fmt.Errorf("read link failed: %v", err)
	}
	return realPath, nil
}

func NullFile() (*os.File, error) {
	return os.OpenFile(constant.NullFilePath, os.O_RDWR, 0666)
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
