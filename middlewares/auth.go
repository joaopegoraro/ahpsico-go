package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joaopegoraro/ahpsico-go/server"
	"github.com/joaopegoraro/ahpsico-go/utils"
)

const (
	TemporaryUserRole = iota
	PatientRole
	DoctorRole
)

const (
	FirstUserRole = PatientRole
	LastUserRole  = DoctorRole
)

type userKey string

var UserKeyCaller = userKey("user")

type AuthUser struct {
	UUID        string
	PhoneNumber string
	Role        int64
	Token       string
}

func Auth(s *server.Server, allowTemporaryUserAccess bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token, claims, err := GetTokenFromRequest(r)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
				RespondAuthError(w, r, s)
				return
			}

			user, err := GetUserDataFromTokenClaims(token, claims)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
				RespondAuthError(w, r, s)
				return
			}

			if !allowTemporaryUserAccess && user.Role == TemporaryUserRole {
				fmt.Printf("Error: Temporary access not allowed")
				RespondAuthError(w, r, s)
				return
			}

			SetTokenHeader(w, user.Token)

			requestContext := context.WithValue(ctx, UserKeyCaller, user)
			next.ServeHTTP(w, r.WithContext(requestContext))
		})
	}
}

func GetTokenFromRequest(r *http.Request) (*jwt.Token, utils.AuthClaims, error) {
	// Get the authorization Token.
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, utils.AuthClaims{}, errors.New("empty auth header")
	}

	// Removes the 'Bearer' prefix of the token
	idTokenSlice := strings.Split(authHeader, " ")
	if len(idTokenSlice) <= 1 {
		return nil, utils.AuthClaims{}, errors.New("invalid auth header")
	}

	headerToken := idTokenSlice[1]
	claims := &utils.AuthClaims{}
	token, err := jwt.ParseWithClaims(headerToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		key := os.Getenv("JWT_SECRET_KEY")
		return []byte(key), nil
	})
	if err != nil || !token.Valid {
		return nil, utils.AuthClaims{}, fmt.Errorf("invalid auth token:\nvalido? %v\n Error: %v\t", token.Valid, err)
	}

	return token, *claims, nil
}

func SetTokenHeader(w http.ResponseWriter, token string) {
	w.Header().Set("token", token)
}

func GetUserDataFromTokenClaims(token *jwt.Token, claims utils.AuthClaims) (AuthUser, error) {
	expirationDate := claims.RegisteredClaims.ExpiresAt.Time

	if expirationDate.Before(time.Now()) {
		return AuthUser{}, errors.New("expired token")
	}

	userToken := token.Raw
	// if the token will be expired within 7 days, it is renewed
	if time.Until(expirationDate).Hours() < float64(time.Hour*24*7) {
		_userToken, err := utils.GenerateJWT(claims.UUID, claims.PhoneNumber, claims.Role)
		if err != nil {
			return AuthUser{}, err
		}
		userToken = _userToken
	}

	return AuthUser{
		UUID:        claims.UUID,
		PhoneNumber: claims.PhoneNumber,
		Role:        claims.Role,
		Token:       userToken,
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
