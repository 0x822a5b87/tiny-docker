package main

import (
	"0x822a5b87/tiny-docker/ns"
	"log"
	"runtime"
)

func main() {
	if runtime.GOOS == "darwin" {
		log.Fatal("error: container engine only runs on Linux (macOS for development only)")
	}

	if err := ns.StartContainer("sh"); err != nil {
		log.Fatal("start container failed:", err)
	}
}
