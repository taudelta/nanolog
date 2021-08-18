package main

import (
	"os"

	log "github.com/taudelta/nanolog"
)

func main() {

	f, err := os.OpenFile("test.log", os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatal().Printf("file error: %v", err)
	}

	defer f.Close()

	log.NoColor()

	log.Init(log.Options{
		Level: log.DebugLevel,
		// File writer overrides default writer
		Debug: log.LoggerOptions{Writer: f},
	})

	log.Debug().Println("debug")

}
