package initializers

import (
	"os"

	"github.com/joaopegoraro/ahpsico-go/server"
	"github.com/twilio/twilio-go"
)

func InitializeAuth(s *server.Server) {
	var TWILIO_ACCOUNT_SID string = os.Getenv("TWILIO_ACCOUNT_SID")
	var TWILIO_AUTH_TOKEN string = os.Getenv("TWILIO_AUTH_TOKEN")
	s.Twilio =  twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: TWILIO_ACCOUNT_SID,
		Password: TWILIO_AUTH_TOKEN,
	})
}
