package daemon

import (
	"github.com/0x822a5b87/tiny-docker/src/handler"
)

func handlePs(request handler.Request) (handler.Response, error) {
	// TODO implement real query.
	return handler.Response{
		Code: 0,
		Msg:  "success",
		Data: []map[string]interface{}{
			{"id": "abc123", "pid": 12345, "status": "running"},
			{"id": "def456", "pid": 67890, "status": "exited"},
		},
	}, nil
}
