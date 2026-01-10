package daemon

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/entity"
	"github.com/sirupsen/logrus"
)

func runContainer(c entity.Container) error {
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

func stopContainer(c entity.Container) error {
	p := getContainerStatusFilePath(c.Id)
	preState, err := readContainerState(p)
	if err != nil {
		logrus.Errorf("error read pre state: %v", err)
		return err
	}

	preState.Status = entity.ContainerExit
	preState.ExitAt = c.ExitAt

	err = writeContainerState(p, preState)
	if err != nil {
		logrus.Errorf("error saving state: %v", err)
		return err
	}

	logrus.Infof("Stop container {%s}", p)
	return nil
}

func getContainerStatusFilePath(id string) string {
	fileRoot := conf.RuntimeDockerdContainerStatus.Get()
	return filepath.Join(fileRoot, id)
}

func readContainerState(p string) (*entity.Container, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		logrus.Errorf("error read pre state : %v", err)
		return nil, err
	}
	preState := &entity.Container{}
	err = json.Unmarshal(data, preState)
	if err != nil {
		logrus.Errorf("error unmarshal pre state : %v", err)
		return nil, err
	}
	return preState, err
}

func writeContainerState(p string, state *entity.Container) error {
	data, err := json.Marshal(state)
	if err != nil {
		logrus.Errorf("error unmarshal pre state : %v", err)
		return err
	}

	if err = os.WriteFile(p, data, 0644); err != nil {
		logrus.Errorf("error saving container: {%v}, {%v}", state, err)
		return err
	}

	return nil
}
