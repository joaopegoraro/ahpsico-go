package main

import (
	"log"
	"net/http"

	"github.com/joaopegoraro/ahpsico-go/initializers"
	"github.com/joaopegoraro/ahpsico-go/server"
)

func main() {
	s := server.NewServer()

	initializers.InitializeContext(s)

	if err := initializers.InitializeDB(s); err != nil {
		log.Fatal(err)
	}

	initializers.InitializeAuth(s)

	initializers.InitializeRoutes(s)

	log.Fatal(http.ListenAndServe(":8000", s.Router))
}
