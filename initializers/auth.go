package initializers

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"github.com/joaopegoraro/ahpsico-go/server"
	"google.golang.org/api/option"
)

func InitializeAuth(s *server.Server) {
	opt := option.WithCredentialsFile("serviceAccount.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	s.Firebase = app
}
