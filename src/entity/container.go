package entity

type ContainerStatus string

var ContainerRunning ContainerStatus = "running"
var ContainerExit ContainerStatus = "exit"

type Container struct {
	Id        string          `json:"id"`
	Pid       int             `json:"pid"`
	Image     string          `json:"image"`
	Command   string          `json:"command"`
	CreatedAt int64           `json:"created_at"`
	ExitAt    int64           `json:"exit_at"`
	Status    ContainerStatus `json:"status"`
	Name      string          `json:"name"`
}

// WaitRequest this request is used to indicate that a detached process is running, and
// requires the daemon process to wait for it.
type WaitRequest struct {
	Id  string `json:"id"`
	Pid int    `json:"pid"`
}
