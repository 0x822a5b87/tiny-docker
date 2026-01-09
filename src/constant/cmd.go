package constant

type Action string

func (action Action) String() string {
	return string(action)
}

const Daemon Action = "daemon"
const InitDaemon Action = "init_daemon"
const InitContainer Action = "init_container"
const Run Action = "run"
const Ps Action = "ps"
const Commit Action = "commit"
