package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joaopegoraro/ahpsico-go/server"
	"github.com/joaopegoraro/ahpsico-go/utils"
)

type userKey string

var UserKeyCaller = userKey("user")

type AuthUser struct {
	UUID        string
	PhoneNumber string
	Role        int64
}

func Auth(s *server.Server) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token, err := GetTokenFromRequest(r)
			if err != nil {
				RespondAuthError(w, r, s)
				return
			}

			user, err := GetUserDataFromToken(token)
			if err != nil {
				RespondAuthError(w, r, s)
				return
			}

			requestContext := context.WithValue(ctx, UserKeyCaller, user)
			next.ServeHTTP(w, r.WithContext(requestContext))
		})
	}
}

func GetTokenFromRequest(r *http.Request) (*jwt.Token, error) {
	// Get the authorization Token.
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("empty auth header")
	}

	// Removes the 'Bearer' prefix of the token
	idTokenSlice := strings.Split(authHeader, " ")
	if len(idTokenSlice) <= 1 {
		return nil, errors.New("invalid auth header")
	}

	headerToken := idTokenSlice[1]
	token, err := jwt.Parse(headerToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			return nil, fmt.Errorf("there's an error with the signing method")
		}

		return os.Getenv("JWT_SECRET_KEY"), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid auth token")
	}

	return token, nil
}

func GetUserDataFromToken(token *jwt.Token) (AuthUser, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return AuthUser{}, errors.New("invalid claims")
	}

	uuid, ok := claims[utils.UuidClaim].(string)
	if !ok {
		return AuthUser{}, errors.New("invalid claims")
	}

	phoneNumber, ok := claims[utils.PhoneNumberClaim].(string)
	if !ok {
		return AuthUser{}, errors.New("invalid claims")
	}

	role, ok := claims[utils.RoleClaim].(int64)
	if !ok {
		return AuthUser{}, errors.New("invalid claims")
	}

	return AuthUser{
		UUID:        uuid,
		PhoneNumber: phoneNumber,
		Role:        role,
	}, nil

}

func GetAuthDataFromContext(ctx context.Context) (AuthUser, uuid.UUID, error) {
	user, ok := ctx.Value(UserKeyCaller).(AuthUser)
	if ok {
		userUuid, err := uuid.FromString(user.UUID)
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
