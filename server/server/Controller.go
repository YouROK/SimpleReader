package server

type Controller interface {
	InitController(serv *Server)
}
