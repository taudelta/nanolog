package main

import (
	log "github.com/stanyx/nanolog"
)

func main() {

	log.Init(log.Options{
		Level: log.DebugLevel,
	})

	log.DEBUG.Println("debug")
	log.INFO.Println("info")
	log.WARN.Println("warn")
	log.ERROR.Println("error")
	log.FATAL.Println("fatal")

	log.Log(log.DebugLevel, "not showed")
}
