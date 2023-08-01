package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/joaopegoraro/ahpsico-go/server"
)

type userKey string

var UserKeyCaller = userKey("user")

type AuthUser struct {
	UID         string
	PhoneNumber string
}

func Auth(s *server.Server) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			idToken, err := GetIdTokenFromRequest(r)
			if err != nil {
				RespondAuthError(w, r, s)
				return
			}

			// Decodes the token.
			decodedToken, err := s.Auth.VerifyIDToken(ctx, idToken)
			if err != nil {
				RespondAuthError(w, r, s)
				return
			}

			// Get the uid from the decoded token, then use it to find and return the user object
			uid := decodedToken.UID
			userRecord, err := s.Auth.GetUser(ctx, uid)
			if err != nil {
				RespondAuthError(w, r, s)
				return
			}

			// pass the user object to the request context
			user := AuthUser{
				UID:         userRecord.UID,
				PhoneNumber: userRecord.PhoneNumber,
			}
			requestContext := context.WithValue(ctx, UserKeyCaller, user)
			next.ServeHTTP(w, r.WithContext(requestContext))
		})
	}
}

func GetIdTokenFromRequest(r *http.Request) (string, error) {
	// Get the authorization Token.
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("empty auth header")
	}

	// Removes the 'Bearer' prefix of the token
	idTokenSlice := strings.Split(authHeader, " ")
	if len(idTokenSlice) <= 1 {
		return "", errors.New("invalid auth header")
	}

	idToken := idTokenSlice[1]
	return idToken, nil
}

func GetAuthDataFromContext(ctx context.Context) (AuthUser, uuid.UUID, error) {
	user, ok := ctx.Value(UserKeyCaller).(AuthUser)
	if ok {
		userUuid, err := uuid.FromString(user.UID)
		if err != nil {
			return user, uuid.Nil, err
		}
		return user, userUuid, nil
	}

	return user, uuid.Nil, errors.New("auth user not found in the context")
}

func RespondAuthError(w http.ResponseWriter, r *http.Request, s *server.Server) {
	s.RespondErrorDetail(w, r, "Invalid auth token", http.StatusUnauthorized)
}
