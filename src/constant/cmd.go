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

const Wait Action = "__wait_request__"
