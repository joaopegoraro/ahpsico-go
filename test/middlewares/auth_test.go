package middlewares_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	firebase_auth "firebase.google.com/go/v4/auth"
	"github.com/gofrs/uuid"
	"github.com/joaopegoraro/ahpsico-go/middlewares"
	"github.com/joaopegoraro/ahpsico-go/server"
)

type mockFirebaseAuth struct {
	Token          *firebase_auth.Token
	TokenError     error
	UserRecord     *firebase_auth.UserRecord
	UserRecorError error
}

func (m *mockFirebaseAuth) VerifyIDToken(ctx context.Context, idToken string) (*firebase_auth.Token, error) {
	return m.Token, m.TokenError
}

func (m *mockFirebaseAuth) GetUser(ctx context.Context, uid string) (*firebase_auth.UserRecord, error) {
	return m.UserRecord, m.UserRecorError
}

func TestGetIdTokenFromRequest(t *testing.T) {
	t.Run("No authorization header returns empty token and error", func(t *testing.T) {
		_, r := createRequestData(t)
		token, err := middlewares.GetIdTokenFromRequest(r)
		if token != "" || err.Error() != "empty auth header" {
			t.Errorf("Token is not empty or err is not expected:\n request: %v \n token: %v \n err: %v", r, token, err.Error())
		}
	})
	t.Run("Token with no Bearer prefix returns empty token and error", func(t *testing.T) {
		_, r := createRequestData(t)
		r.Header.Add("Authorization", "sometoken")
		token, err := middlewares.GetIdTokenFromRequest(r)
		if token != "" || err.Error() != "invalid auth header" {
			t.Errorf("Token is not empty or err is not expected:\n request: %v \n token: %v \n err: %v", r, token, err.Error())
		}
	})
	t.Run("Valid auth header returns token and no error", func(t *testing.T) {
		_, r := createRequestData(t)
		expectedToken := "sometoken"
		r.Header.Add("Authorization", "Bearer "+expectedToken)
		token, err := middlewares.GetIdTokenFromRequest(r)
		if token != expectedToken || err != nil {
			t.Errorf("Token is not empty or err is not expected:\n request: %v \n token: %v \n err: %v", r, expectedToken, err.Error())
		}
	})
}

func TestGetAuthDataFromContext(t *testing.T) {
	userUuid, _ := uuid.NewV1()
	validUserWithBadUUID := middlewares.AuthUser{}
	validUser := middlewares.AuthUser{UID: userUuid.String()}

	tests := []struct {
		name    string
		ctx     context.Context
		want    middlewares.AuthUser
		want1   uuid.UUID
		wantErr bool
	}{
		{
			name:    "user not found returns no user, nil uuid and valid error",
			ctx:     context.Background(),
			want:    middlewares.AuthUser{},
			want1:   uuid.Nil,
			wantErr: true,
		},
		{
			name:    "user found with invalid id returns user, nil uuid and valid error",
			ctx:     context.WithValue(context.Background(), middlewares.UserKeyCaller, validUserWithBadUUID),
			want:    validUserWithBadUUID,
			want1:   uuid.Nil,
			wantErr: true,
		},
		{
			name:    "user found with valid id returns user, uuid and nil error",
			ctx:     context.WithValue(context.Background(), middlewares.UserKeyCaller, validUser),
			want:    validUser,
			want1:   userUuid,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := middlewares.GetAuthDataFromContext(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAuthDataFromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAuthDataFromContext() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetAuthDataFromContext() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRespondAuthError(t *testing.T) {
	s := server.NewServer()
	t.Run("When called s.RespondErrorDetail with error detail 401 status", func(t *testing.T) {
		w, r := createRequestData(t)
		middlewares.RespondAuthError(w, r, s)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("code not 401: %v", w)
		}
	})
}

func TestAuth(t *testing.T) {
	s := server.NewServer()
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	t.Run("Invalid token header responds with error", func(t *testing.T) {
		w, r := createRequestData(t)
		middlewares.Auth(s)(emptyHandler).ServeHTTP(w, r)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("code not 401: %v", w)
		}
	})
	t.Run("Error verifying token responds with error", func(t *testing.T) {
		s.Auth = &mockFirebaseAuth{TokenError: errors.New("some error")}
		w, r := createRequestData(t)
		r.Header.Add("Authorization", "Bearer some token")
		middlewares.Auth(s)(emptyHandler).ServeHTTP(w, r)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("code not 401: %v", w)
		}
	})
	t.Run("Error retrieving user data responds with error", func(t *testing.T) {
		s.Auth = &mockFirebaseAuth{
			Token:          &firebase_auth.Token{UID: "some uid"},
			UserRecorError: errors.New("some error"),
		}
		w, r := createRequestData(t)
		r.Header.Add("Authorization", "Bearer some token")
		middlewares.Auth(s)(emptyHandler).ServeHTTP(w, r)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("code not 401: %v", w)
		}
	})
	t.Run("Successfuly retrieving user data responds with error", func(t *testing.T) {
		expectedUID, _ := uuid.NewV1()
		expectedPhoneNumber := "some phone number"
		expectedUser := &firebase_auth.UserRecord{UserInfo: &firebase_auth.UserInfo{
			UID:         expectedUID.String(),
			PhoneNumber: expectedPhoneNumber,
		}}
		s.Auth = &mockFirebaseAuth{
			Token:      &firebase_auth.Token{UID: "some uid"},
			UserRecord: expectedUser,
		}
		w, r := createRequestData(t)
		r.Header.Add("Authorization", "Bearer some token")
		middlewares.Auth(s)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUser, _, err := middlewares.GetAuthDataFromContext(r.Context())
			if err != nil {
				t.Fatal(err)
			}
			w.WriteHeader(http.StatusAccepted)
			if authUser.UID != expectedUID.String() || authUser.PhoneNumber != expectedPhoneNumber {
				t.Errorf("auth user with unexpected info: %v", authUser)
			}
		})).ServeHTTP(w, r)
		if w.Code != http.StatusAccepted {
			t.Errorf("code not 202: %v", w)
		}
	})
}

func createRequestData(t *testing.T) (*httptest.ResponseRecorder, *http.Request) {
	req, err := http.NewRequest("GET", "/some-endpoint", nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	return rr, req
}
