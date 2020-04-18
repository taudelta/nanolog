package main

import (
	"os"

	log "github.com/stanyx/nanolog"
)

func main() {

	f, err := os.OpenFile("test.log", os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.FATAL.Printf("file error: %v", err)
	}

	defer f.Close()

	log.NoColor()

	log.Init(log.Options{
		Level: log.DebugLevel,
		// File writer overrides default writer
		Debug: f,
	})

	log.DEBUG.Println("debug")

}
