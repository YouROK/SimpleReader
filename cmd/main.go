package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-martini/martini"

	"SimpleReader/server"
	"SimpleReader/server/models/global"
	"SimpleReader/server/models/storage"
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
	host := ""
	port := "8092"
	chdir := "/home/yourok/sr"
	//TODO remove
	//chdir = "/Users/yourok/Projects/GO/SimpleReader"
	if len(os.Args) >= 2 {
		host = os.Args[1]
	}

	if len(os.Args) >= 3 {
		port = os.Args[2]
	}

	if len(os.Args) >= 4 {
		chdir = os.Args[3]
	}

	os.Chdir(chdir)
	martini.Root = chdir
	log.Println("SR run with:", host, port, chdir)
	go server.Start(host, port)
	catchExitSignals()
	global.Stoping = true
	storage.GetStorage().Exit()
	log.Println("Exit")
}
