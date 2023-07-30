package initializers

import (
	"database/sql"

	"github.com/joaopegoraro/ahpsico-go/database/db"
	"github.com/joaopegoraro/ahpsico-go/server"

	_ "github.com/mattn/go-sqlite3"
)

func InitializeDB(s *server.Server) error {
	sqldb, err := sql.Open("sqlite3", "database/db.sqlite3")
	if err != nil {
		return err
	}

	if s != nil {
		s.Queries = db.New(sqldb)
	}

	return nil
}
