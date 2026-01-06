package util

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

func UnTar(tarPath string, destDir string) error {
	if _, err := exec.Command("tar", "-xvf", tarPath, "-C", destDir).CombinedOutput(); err != nil {
		logrus.Errorf("untar error: {%s}", err.Error())
		return err
	}
	return nil
}
