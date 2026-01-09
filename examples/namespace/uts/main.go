package main

import (
	"log"
	"runtime"
)

func main() {
	if runtime.GOOS == "darwin" {
		log.Fatal("error: daemon engine only runs on Linux (macOS for development only)")
	}

	if err := StartContainer("sh"); err != nil {
		log.Fatal("start daemon failed:", err)
	}
}
