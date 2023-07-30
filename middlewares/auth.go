package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/joaopegoraro/ahpsico-go/server"
)

type userKey string

var userKeyCaller = userKey("user")

type AuthUser struct {
	UID         string
	PhoneNumber string
}

func Auth(s *server.Server) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the authorization Token.
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handleAuthError(w, r, s)
				return
			}

			// Removes the 'Bearer' prefix of the token
			idTokenSlice := strings.Split(authHeader, " ")
			idToken := idTokenSlice[1]

			// Get the firebase auth client
			auth, err := s.Firebase.Auth(s.Ctx)
			if err != nil {
				handleAuthError(w, r, s)
				return
			}

			// Decodes the token.
			decodedToken, err := auth.VerifyIDToken(s.Ctx, idToken)
			if err != nil {
				handleAuthError(w, r, s)
				return
			}

			// Get the uid from the decoded token, then use it to find and return the user object
			uid := decodedToken.UID
			userRecord, err := auth.GetUser(s.Ctx, uid)
			if err != nil {
				handleAuthError(w, r, s)
				return
			}

			// pass the user object to the request context
			user := AuthUser{
				UID:         userRecord.UID,
				PhoneNumber: userRecord.PhoneNumber,
			}
			requestContext := context.WithValue(r.Context(), userKeyCaller, user)
			next.ServeHTTP(w, r.WithContext(requestContext))
		})
	}
}

func GetAuthUserFromContext(ctx context.Context) (AuthUser, error) {
	user, ok := ctx.Value(userKeyCaller).(AuthUser)
	if ok {
		return user, nil
	}
	return user, errors.New("auth user not found in the context")
}

func handleAuthError(w http.ResponseWriter, r *http.Request, s *server.Server) {
	s.RespondError(w, r, "Invalid auth token", http.StatusUnauthorized)
}
