package initializers

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"github.com/joaopegoraro/ahpsico-go/server"
)

func InitializeAuth(s *server.Server) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	s.Firebase = app
}
