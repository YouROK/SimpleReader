package main

import (
	"SimpleReader/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func catchExitSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	log.Println("Catch signals")
	for sig := range c {
		log.Println("Signal to exit:", sig)
		return
	}
}

func main() {
	go server.Start()
	catchExitSignals()
}
