package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/joaopegoraro/ahpsico-go/server"
)

func TestRespond(t *testing.T) {
	s := server.NewServer()
	t.Run("Nil data only writes provided status", func(t *testing.T) {
		w, r := createRequestData(t)
		s.Respond(w, r, nil, http.StatusOK)
		if w.Body.String() != "" || w.Code != http.StatusOK {
			t.Errorf("Body not nil or code not 200: %v", w)
		}
	})
	t.Run("Error marshaling data returns 500 with error body", func(t *testing.T) {
		jsonErr := json.UnsupportedTypeError{Type: reflect.ValueOf(func() {}).Type()}
		expectedBody, err := json.Marshal(server.Error{Detail: jsonErr.Error()})
		if err != nil {
			t.Fatal(err)
		}
		w, r := createRequestData(t)
		s.Respond(w, r, func() {}, http.StatusOK)
		if w.Body.String() != string(expectedBody) || w.Code != http.StatusInternalServerError {
			t.Errorf("Body not expected or code not 500: %v", w)
		}
	})
	t.Run("Success marshaling data returns provided body and status", func(t *testing.T) {
		data := server.Error{Detail: "Hello"}
		expectedBody, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}
		w, r := createRequestData(t)
		s.Respond(w, r, data, http.StatusCreated)
		if w.Body.String() != string(expectedBody) || w.Code != http.StatusCreated {
			t.Errorf("Body or code not expected: %v", w)
		}
	})
}

func TestRespondNoContent(t *testing.T) {
	s := server.NewServer()
	t.Run("s.Respond is called with nil data and http.StatusNoContent", func(t *testing.T) {
		w, r := createRequestData(t)
		s.RespondNoContent(w, r)
		if w.Body.String() != "" || w.Code != http.StatusNoContent {
			t.Errorf("Body not nil or code not 204: %v", w)
		}
	})
}

func TestRespondOk(t *testing.T) {
	s := server.NewServer()
	t.Run("Nil data calls s.RespondNoContent", func(t *testing.T) {
		w, r := createRequestData(t)
		s.RespondOk(w, r, nil)
		if w.Body.String() != "" || w.Code != http.StatusNoContent {
			t.Errorf("Body not nil or code not 204: %v", w)
		}
	})
	t.Run("Responds provided body and StatusOK", func(t *testing.T) {
		data := server.Error{Detail: "Hello"}
		expectedBody, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}
		w, r := createRequestData(t)
		s.RespondOk(w, r, data)
		if w.Body.String() != string(expectedBody) || w.Code != http.StatusOK {
			t.Errorf("Body not expected or code not 200: %v", w)
		}
	})
}

func TestRespondError(t *testing.T) {
	s := server.NewServer()
	t.Run("Status less than 1 calls s.Respond with provided error body and 500 status", func(t *testing.T) {
		errorData := server.Error{Detail: "Hello", Status: 0}
		expectedBody, err := json.Marshal(errorData)
		if err != nil {
			t.Fatal(err)
		}
		w, r := createRequestData(t)
		s.RespondError(w, r, errorData)
		if w.Body.String() != string(expectedBody) || w.Code != http.StatusInternalServerError {
			t.Errorf("Body not expected or code not 500: %v", w)
		}
	})
	t.Run("Status more or equal to 1 calls s.Respond with provided error and status", func(t *testing.T) {
		errorData := server.Error{Detail: "Hello", Status: http.StatusBadRequest}
		expectedBody, err := json.Marshal(errorData)
		if err != nil {
			t.Fatal(err)
		}
		w, r := createRequestData(t)
		s.RespondError(w, r, errorData)
		if w.Body.String() != string(expectedBody) || w.Code != http.StatusBadRequest {
			t.Errorf("Body not expected or code not 400: %v", w)
		}
	})
}

func TestRespondErrorStatus(t *testing.T) {
	s := server.NewServer()
	t.Run("Status less than 1 calls s.Respond with nil data and 500 status", func(t *testing.T) {
		w, r := createRequestData(t)
		s.RespondErrorStatus(w, r, 0)
		if w.Body.String() != "" || w.Code != http.StatusInternalServerError {
			t.Errorf("Body not empty or code not 500: %v", w)
		}
	})
	t.Run("Status more or equal to 1 calls s.Respond with nil data and provided status", func(t *testing.T) {
		w, r := createRequestData(t)
		s.RespondErrorStatus(w, r, http.StatusBadRequest)
		if w.Body.String() != "" || w.Code != http.StatusBadRequest {
			t.Errorf("Body not empty or code not 400: %v", w)
		}
	})
}

func TestRespondErrorDetail(t *testing.T) {
	s := server.NewServer()
	t.Run("Blank detail calls s.RespondErrorStatus with provided status", func(t *testing.T) {
		w, r := createRequestData(t)
		s.RespondErrorDetail(w, r, " ", http.StatusBadRequest)
		if w.Body.String() != "" || w.Code != http.StatusBadRequest {
			t.Errorf("Body not empty or code not 400: %v", w)
		}
	})
	t.Run("Non blank detail calls s.RespondError with error data containing provided detail and status", func(t *testing.T) {
		errorData := server.Error{Detail: "Hello", Status: http.StatusBadRequest}
		expectedBody, err := json.Marshal(errorData)
		if err != nil {
			t.Fatal(err)
		}
		w, r := createRequestData(t)
		s.RespondErrorDetail(w, r, errorData.Detail, errorData.Status)
		if w.Body.String() != string(expectedBody) || w.Code != http.StatusBadRequest {
			t.Errorf("Body not expected or code not 400: %v", w)
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
