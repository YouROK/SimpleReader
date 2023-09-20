package server

import (
	"log"
	"net/http"
	"reflect"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"

	"SimpleReader/server_old/models/storage"
)

type Server struct {
	martini *martini.ClassicMartini
	Host    string
}

func CreateServer(host string) *Server {
	s := storage.GetStorage()
	s.SetPath("db/storage")

	server := Server{}
	r := martini.NewRouter()
	m := martini.New()
	//m.Use(martini.Logger())
	m.Use(gzip.All())
	m.Use(martini.Recovery())
	m.Use(martini.Static("public"))
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)

	server.martini = &martini.ClassicMartini{m, r}
	server.martini.Use(storage.Sessions())
	server.martini.Use(martini.Static("db/storage/"))
	views := "views"
	//if martini.Env == martini.Dev {
	//	views = "server/" + views
	//}
	server.martini.Use(render.Renderer(render.Options{Directory: views, Extensions: []string{".go.html", ".html"}}))
	server.martini.MapTo(s, (*storage.Storages)(nil))
	server.Host = host
	return &server
}

func (s *Server) GetMartini() *martini.ClassicMartini {
	return s.martini
}

func (s *Server) AddControler(controler interface{}) {
	if c, ok := controler.(Controller); ok {
		c.InitController(s)
	} else {
		panic(reflect.TypeOf(controler).String() + " not implement Controller interface")
	}
}

func (s *Server) Run() {
	log.Fatalln(http.ListenAndServe(s.Host, s.martini))
}
