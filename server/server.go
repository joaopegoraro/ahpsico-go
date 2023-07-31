package server

import (
	"encoding/json"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/joaopegoraro/ahpsico-go/database/db"
)

type Server struct {
	Router   *chi.Mux
	Queries  *db.Queries
	Firebase *firebase.App
}

type Error struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	Status   int    `json:"status,omitempty"`

	//	"type": "https://example.com/probs/out-of-credit",
	//	"title": "You do not have enough credit.",
	//	"detail": "Your current balance is 30, but that costs 50.",
	//	"instance": "/account/12345/msgs/abc",
	//	"status": 400,
}

const Success = "SUCCESS"

func NewServer() *Server {
	s := &Server{}
	return s
}

func (s *Server) Respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	if data == nil {
		w.WriteHeader(status)
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		errData, err := json.Marshal(Error{Detail: err.Error()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errData)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
}

func (s *Server) RespondOk(w http.ResponseWriter, r *http.Request, data interface{}) {
	if data == nil {
		s.RespondNoContent(w, r)
		return
	}

	s.Respond(w, r, data, http.StatusOK)
}

func (s *Server) RespondNoContent(w http.ResponseWriter, r *http.Request) {
	s.Respond(w, r, nil, http.StatusNoContent)
}

func (s *Server) RespondError(w http.ResponseWriter, r *http.Request, err Error) {
	if err.Status < 1 {
		s.Respond(w, r, err, http.StatusInternalServerError)
		return
	}
	s.Respond(w, r, err, err.Status)
}

func (s *Server) RespondErrorStatus(w http.ResponseWriter, r *http.Request, status int) {
	if status < 1 {
		s.Respond(w, r, nil, http.StatusInternalServerError)
		return
	}
	s.Respond(w, r, nil, status)
}

func (s *Server) RespondErrorDetail(w http.ResponseWriter, r *http.Request, detail string, status int) {
	if strings.TrimSpace(detail) == "" {
		s.Respond(w, r, nil, status)
		return
	}

	s.Respond(w, r, Error{Detail: detail}, status)
}

func (s *Server) Decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
