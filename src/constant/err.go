package constant

import (
	"fmt"
)

type Err struct {
	ErrorCode int
	ErrorText string
}

func (e Err) Error() string {
	return fmt.Sprintf("code = [%d], error = [%s]", e.ErrorCode, e.ErrorText)
}

func (e Err) Wrap(wrapErr error) error {
	return fmt.Errorf(e.Error(), wrapErr.Error())
}

func (e Err) WrapMessage(msg string) error {
	return fmt.Errorf(e.Error(), msg)
}

var (
	ErrMalformedType                = Err{ErrorCode: 100000, ErrorText: "Malformed type"}
	ErrProcsEmpty                   = Err{ErrorCode: 100001, ErrorText: "Procs is nil or empty"}
	ErrProcsPidNotFound             = Err{ErrorCode: 100002, ErrorText: "Pid not found"}
	ErrCreateCgroup                 = Err{ErrorCode: 100003, ErrorText: "Error create cgroup"}
	ErrMalformedArgs                = Err{ErrorCode: 100004, ErrorText: "Malformed args"}
	ErrMountRootFS                  = Err{ErrorCode: 100005, ErrorText: "Mount rootfs to itself error: %v"}
	ErrCreateUdsServer              = Err{ErrorCode: 100006, ErrorText: "Create uds server error: %v"}
	ErrMalformedUdsReq              = Err{ErrorCode: 100007, ErrorText: "Malformed uds request error: %v"}
	ErrMalformedUdsRsp              = Err{ErrorCode: 100008, ErrorText: "Malformed uds response error: %v"}
	ErrUnsupportedAction            = Err{ErrorCode: 100009, ErrorText: "Unsupported action error: %v"}
	ErrIllegalUdsServerStatus       = Err{ErrorCode: 100010, ErrorText: "Illegal UDS server status: %v"}
	ErrProcessTerminalAndDaemonMode = Err{ErrorCode: 100011, ErrorText: "Interactive (-it) and detach (-d) modes are mutually exclusive"}
	ErrExecCommand                  = Err{ErrorCode: 100012, ErrorText: "Execute command error: %v"}
	ErrMalformedLogsArgs            = Err{ErrorCode: 100013, ErrorText: "Mini-docker requires 1 argument"}
	ErrResourceNotFound             = Err{ErrorCode: 100014, ErrorText: "resource not found"}
	ErrResourceExists               = Err{ErrorCode: 100015, ErrorText: "resource exists"}
	ErrResourcePoolIsEmpty          = Err{ErrorCode: 100016, ErrorText: "resource Pool is empty, cannot put resource back"}
	ErrOutOfRange                   = Err{ErrorCode: 100017, ErrorText: "out of range"}
	ErrInvalidPos                   = Err{ErrorCode: 100018, ErrorText: "invalid position (out of range)"}
	ErrInvalidSize                  = Err{ErrorCode: 100019, ErrorText: "size must be greater than 0"}
	ErrNetworkVersion               = Err{ErrorCode: 100020, ErrorText: "only IPv4 is supported"}
	ErrResourceNotExists            = Err{ErrorCode: 100021, ErrorText: "resource not exists"}
	ErrDeviceIsBusy                 = Err{ErrorCode: 100022, ErrorText: "device is busy"}
)
