package main

import (
	"SimpleReader"
	"SimpleReader/models/global"
	"SimpleReader/models/storage"
	"github.com/go-martini/martini"
	"log"
	"os"
	"os/signal"
	"runtime"
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
	runtime.GOMAXPROCS(4)
	host := ""
	port := "9000"
	chdir := "/home/yourok/sr" //os.Getenv("OPENSHIFT_DATA_DIR")

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
	go SimpleReader.Start(host, port)
	catchExitSignals()
	global.Stoping = true
	storage.GetStorage().Exit()
	log.Println("Exit")
}
