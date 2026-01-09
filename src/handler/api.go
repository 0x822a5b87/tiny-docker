package handler

import (
	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	log "github.com/sirupsen/logrus"
)

func SendPsRequest() error {
	req, err := ParamsIntoRequest[any](constant.Commit, "")
	err, rsp := sendRequest(req)
	if err != nil {
		log.Errorf("error sending ps request: %v\n", err)
		return err
	}
	log.Infof("rsp: %v", rsp)
	return nil

}

func SendCommitRequest(commands conf.CommitCommands) error {
	req, err := ParamsIntoRequest[conf.CommitCommands](constant.Commit, commands)
	if err != nil {
		return err
	}
	err, rsp := sendRequest(req)
	if err != nil {
		log.Errorf("error sending commit request: %v\n", err)
		return err
	}
	log.Infof("rsp: %v", rsp)
	return nil

}
