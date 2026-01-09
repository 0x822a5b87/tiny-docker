package handler

import (
	"encoding/json"
	"net"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/sirupsen/logrus"
)

func sendRequest(req *Request) (error, Response) {
	if !isUdsServerRunning() {
		logrus.Fatalf("UDS server is not running")
		return constant.ErrIllegalUdsServerStatus, Response{}
	}

	// connect client
	conn, err := net.DialUnix(constant.OS, nil, &net.UnixAddr{
		Name: conf.RuntimeDockerdUdsFile.Get(),
		Net:  constant.OS,
	})
	if err != nil {
		logrus.Errorf("failed to dial mini-dockerd: %v\n", err.Error())
		return err, Response{}
	}
	defer conn.Close()

	reqData, _ := json.Marshal(req)
	_, err = conn.Write(reqData)
	if err != nil {
		logrus.Errorf("failed to send uds request: %v\n", err.Error())
		return err, Response{}
	}

	buf := make([]byte, 1024*10)
	n, err := conn.Read(buf)
	if err != nil {
		logrus.Errorf("failed to read uds response: %v\n", err)
		return err, Response{}
	}

	var rsp Response
	err = json.Unmarshal(buf[:n], &rsp)
	if err != nil {
		logrus.Errorf("error unmarshal uds response: %v\n", err.Error())
		return err, Response{}
	}

	return nil, rsp
}
