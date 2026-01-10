package daemon

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/entity"
	"github.com/sirupsen/logrus"
)

func saveContainer(c entity.Container) error {
	data, err := json.Marshal(c)
	if err != nil {
		logrus.Error("error serialize container: {%s}, {%v}", c, err)
		return err
	}
	p := getContainerStatusFilePath(c.Id)
	if err = os.WriteFile(p, data, 0644); err != nil {
		logrus.Errorf("error saving container: {%v}, {%v}", c, err)
		return err
	}
	logrus.Infof("Saving container in file {%s}", p)
	return nil
}

func getContainerStatusFilePath(id string) string {
	fileRoot := conf.RuntimeDockerdContainerStatus.Get()
	return filepath.Join(fileRoot, id)
}
