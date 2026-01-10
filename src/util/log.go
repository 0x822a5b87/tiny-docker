package util

import (
	"fmt"
)

func LogWithoutExtraInfo(msg any) {
	fmt.Printf("%v", msg)
}
