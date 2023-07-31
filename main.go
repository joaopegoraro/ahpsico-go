package main

import (
	"log"
	"net/http"

	"github.com/joaopegoraro/ahpsico-go/initializers"
	"github.com/joaopegoraro/ahpsico-go/server"
)

func main() {
	s := server.NewServer()

	if err := initializers.InitializeDB(s); err != nil {
		log.Fatal(err)
	}

	initializers.InitializeAuth(s)

	initializers.InitializeRoutes(s)

	log.Print("Serving on :8000")

	log.Fatal(http.ListenAndServe(":8000", s.Router))
}
