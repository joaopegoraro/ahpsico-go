package initializers

import (
	"database/sql"
	"log"

	"github.com/joaopegoraro/ahpsico-go/database/db"
	"github.com/joaopegoraro/ahpsico-go/server"

	_ "github.com/mattn/go-sqlite3"
)

func InitializeDB(s *server.Server) {
	sqldb, err := sql.Open("sqlite3", "database/db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	if s == nil {
		log.Fatal()
	}

	s.Queries = db.New(sqldb)
}
