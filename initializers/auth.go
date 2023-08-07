package initializers

import (
	"context"
	"log"

	"github.com/joaopegoraro/ahpsico-go/server"
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
