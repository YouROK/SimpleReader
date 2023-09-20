package server

import (
	"SimpleReader/server_old/controllers"
	"SimpleReader/server_old/server"
)

type SimpleReader struct {
	server *server.Server
}

func Start(host, port string) {
	reader := SimpleReader{}

	reader.server = server.CreateServer(host + ":" + port)

	//Routes
	//root
	reader.server.AddControler(&controllers.MainPage{})
	//auth
	reader.server.AddControler(&controllers.Authorization{})
	//books
	reader.server.AddControler(&controllers.Book{})

	reader.server.AddControler(&controllers.Styles{})

	reader.server.AddControler(&controllers.Upload{})

	reader.server.AddControler(&controllers.Library{})

	reader.server.AddControler(&controllers.Hello{}) // test

	reader.server.Run()
}
