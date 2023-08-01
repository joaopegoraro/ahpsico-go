package initializers

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"github.com/joaopegoraro/ahpsico-go/server"
	"google.golang.org/api/option"
)

func InitializeAuth(s *server.Server) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing auth client: %v\n", err)
		return
	}

	s.Firebase = app
	s.Auth = auth
}
