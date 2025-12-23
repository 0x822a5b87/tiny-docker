package constant

import "fmt"

type Err struct {
	ErrorCode int
	ErrorText string
}

func (e Err) Error() string {
	return fmt.Sprintf("code = [%d], error = [%s]", e.ErrorCode, e.ErrorText)
}

var (
	ErrMalformedType    = Err{ErrorCode: 100000, ErrorText: "malformed type"}
	ErrProcsEmpty       = Err{ErrorCode: 100001, ErrorText: "procs is nil or empty"}
	ErrProcsPidNotFound = Err{ErrorCode: 100002, ErrorText: "pid not found"}
	ErrCreateCgroup     = Err{ErrorCode: 100003, ErrorText: "error create cgroup"}
)
