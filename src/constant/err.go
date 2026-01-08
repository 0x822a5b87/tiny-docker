package constant

import "fmt"

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

var (
	ErrMalformedType    = Err{ErrorCode: 100000, ErrorText: "malformed type"}
	ErrProcsEmpty       = Err{ErrorCode: 100001, ErrorText: "procs is nil or empty"}
	ErrProcsPidNotFound = Err{ErrorCode: 100002, ErrorText: "pid not found"}
	ErrCreateCgroup     = Err{ErrorCode: 100003, ErrorText: "error create cgroup"}
	ErrMalformedArgs    = Err{ErrorCode: 100004, ErrorText: "malformed args"}
	ErrMountRootFS      = Err{ErrorCode: 100005, ErrorText: "Mount rootfs to itself error: %v"}
)
