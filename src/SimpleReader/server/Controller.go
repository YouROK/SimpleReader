package server

import ()

type Controller interface {
	InitController(serv *Server)
}
