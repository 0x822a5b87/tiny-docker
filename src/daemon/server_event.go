package daemon

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"syscall"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/entity"
	"github.com/0x822a5b87/tiny-docker/src/util"
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

func stopContainers(containers []entity.Container) error {
	// TODO it must return a error list
	var err error
	for _, container := range containers {
		err = stopContainer(container)
		if err != nil {
			logrus.Errorf("error stop container: %s", container.Id)
		}
	}
	return err
}

func stopContainer(c entity.Container) error {
	p := getContainerStatusFilePath(c.Id)
	preState, err := readContainerState(p)
	if err != nil {
		logrus.Errorf("error read pre state: %v", err)
		return err
	}

	if preState.Status == entity.ContainerRunning {
		if err = util.KillProcessByPID(preState.Pid, 9); err != nil {
			if errors.Is(err, syscall.ESRCH) {
				logrus.Debugf("process %d does not exist, ignore error", preState.Pid)
			} else {
				return err
			}
		}
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

func logs(containerId string) (string, error) {
	logFile := getContainerLogFilePath(containerId)
	data, err := os.ReadFile(logFile)
	if err != nil {
		logrus.Errorf("error read log file: %v", err)
		return "", err
	}
	return string(data), nil
}

func ps(command conf.PsCommand) ([]entity.Container, error) {
	allContainers, err := readAllContainers()
	if err != nil {
		logrus.Errorf("error reading all containers: %v", err)
		return nil, err
	}

	if command.All {
		return allContainers, nil
	}

	targetContainers := make([]entity.Container, 0)
	for _, container := range allContainers {
		if container.Status != entity.ContainerRunning {
			continue
		}
		targetContainers = append(targetContainers, container)
	}
	return targetContainers, nil
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

func readAllContainers() ([]entity.Container, error) {
	p := conf.RuntimeDockerdContainerStatus.Get()
	containerData, err := util.ReadAllFilesInDir(p)
	if err != nil {
		logrus.Errorf("error read all containers : %v", err)
		return nil, err
	}

	containers := make([]entity.Container, 0)
	for _, data := range containerData {
		var container entity.Container
		if err = json.Unmarshal(data, &container); err != nil {
			logrus.Errorf("error unmarshal container : %v", err)
			continue
		}
		containers = append(containers, container)
	}
	return containers, nil
}
