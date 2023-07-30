package initializers

import (
	"context"

	"github.com/joaopegoraro/ahpsico-go/server"
)

func InitializeContext(s *server.Server) {
	s.Ctx = context.Background()
}
