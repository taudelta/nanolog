package main

import (
	log "github.com/taudelta/nanolog"
)

func main() {

	log.Init(log.Options{
		Level: log.DebugLevel,
		Debug: log.LoggerOptions{
			Prefix: "%v> ",
			Flags:  log.LstdFlags | log.Llongfile,
		},
	})

	log.Debug().Println("debug")
	log.Info().Println("info")
	log.Warn().Println("warn")
	log.Error().Println("error")
	log.Fatal().Println("fatal")

	log.Log(log.DebugLevel, "not showed")
}
