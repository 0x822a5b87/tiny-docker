package constant

type Action string

func (action Action) String() string {
	return string(action)
}

const Daemon Action = "daemon"
const InitDaemon Action = "init_daemon"
const InitContainer Action = "init_container"
const Run Action = "run"
const Stop Action = "stop"
const Ps Action = "ps"
const Commit Action = "commit"
const Logs Action = "logs"
const Exec Action = "exec"
const Network Action = "network"
const NetworkCreate Action = "create"
const NetworkConnect Action = "connect"
const NetworkRm Action = "rm"
const NetworkInspect Action = "inspect"

const Wait Action = "__wait_request__"
