package handler

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/sirupsen/logrus"
)

func CreateUdsServer() error {
	if isUdsServerRunning() {
		logrus.Fatalf("UDS server is already running")
		return nil
	}

	udsPath := conf.RuntimeDockerdUdsFile.Get()

	err := util.EnsureFileExists(udsPath)
	if err != nil {
		logrus.Errorf("UDS server uds path does not exist: {%s}", udsPath)
		return err
	}
	_ = os.Remove(udsPath)

	network := constant.OS
	listener, err := net.ListenUnix(network, &net.UnixAddr{
		Name: udsPath,
		Net:  network,
	})

	if err != nil {
		return constant.ErrCreateUdsServer.Wrap(err)
	}

	defer listener.Close()

	if err = os.Chmod(udsPath, 0600); err != nil {
		return err
	}

	pidFile := conf.RuntimeDockerdUdsPidFile.Get()
	if err = util.EnsureFileExists(pidFile); err != nil {
		logrus.Errorf("UDS pid file does not exist: {%s}", pidFile)
		return err
	}

	{
		// TODO delete me
		file, err := os.Open(pidFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	if err = os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0600); err != nil {
		logrus.Errorf("error write file : %v", pidFile)
		return err
	}

	logrus.Infof("start UDS for mini-dockerd on : %s", udsPath)

	for {
		conn, err := listener.AcceptUnix()
		if err != nil {
			logrus.Errorf("error accept unix：%v\n", err)
			continue
		}
		go handleClient(conn)
	}
}

func isUdsServerRunning() bool {
	pidFile := conf.RuntimeDockerdUdsPidFile.Get()
	if err := util.EnsureFileExists(pidFile); err != nil {
		logrus.Errorf("error check file status : %v", err)
		return false
	}
	pidStr, err := os.ReadFile(pidFile)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		logrus.Errorf("error read UDS pid file: %v", err)
		return false
	}
	pid := 0
	_, err = fmt.Sscanf(string(pidStr), "%d", &pid)
	if err != nil {
		return false
	}

	// TODO this is unsafe because the process can be a irrelative process
	// check if mini-dockerd is alive by sending signal 0
	if err = syscall.Kill(pid, 0); err != nil {
		return false
	}

	return true
}

func handleClient(conn *net.UnixConn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		logrus.Errorf("error read client data：%v\n", err)
		return
	}

	var req Request
	err = json.Unmarshal(buf[:n], &req)
	if err != nil {
		resp, _ := ErrorResponse(err, constant.ErrMalformedUdsReq)
		if err = sendResponse(conn, resp); err != nil {
			logrus.Error("error send response: %s\n", err.Error())
		}
		return
	}

	handleRsp(conn, req)
}

func handleRsp(conn *net.UnixConn, req Request) {
	rsp, err := handleRequest(req)
	if err != nil {
		logrus.Errorf("handle request error: %s\n", err.Error())
		return
	}
	if err = sendResponse(conn, rsp); err != nil {
		logrus.Error("error send response: %s\n", err.Error())
	}
}

func sendResponse(conn *net.UnixConn, resp Response) error {
	respData, _ := json.Marshal(resp)
	_, err := conn.Write(respData)
	return err
}
