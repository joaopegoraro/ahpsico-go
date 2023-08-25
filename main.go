package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joaopegoraro/ahpsico-go/initializers"
	"github.com/joaopegoraro/ahpsico-go/server"
)

func main() {
	s := server.NewServer()

	initializers.InitializeEnv()
	initializers.InitializeDB(s)
	initializers.InitializeAuth(s)
	initializers.InitializeRoutes(s)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal()
	}

	log.Printf("Serving on :%s", port)
	log.Fatal(http.ListenAndServe(port, s.Router))
}
