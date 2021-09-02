package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRegisterSuccessHandler tests that Register responds successfully
// given valid signup information
func TestRegisterSuccessHandler(t *testing.T) {

	var jsonStr = []byte(`{"name":"Aaron von Kreisler", "email":"test@test.com", "password":"password"}`)
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))

	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record responses
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Register)

	// Call ServeHTTP directly and pass in our Request and Response Recorder
	handler.ServeHTTP(rr, req)

	// Check that the status code is what we expect

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var m map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &m)

	if err != nil {
		t.Error("Failed unmarshaling response")
	}

	if m["email"] != "test@test.com" {
		t.Errorf("Failed email: got %s want %s", m["email"], "test@test.com")
	}

	if m["password"] != nil {
		t.Error("Received password field which should not be included. Convert internal user object to UserResponse object")
	}

	if m["id"] == "" {
		t.Error("Missing UUID")
	}

}

// TestRegisterFailWithBadEmail tests whether Register hangler returns the correct
// status code / error message when attempting to register with an invalid email
func TestRegisterFailWithBadEmail(t *testing.T) {
	var jsonStr = []byte(`"name":"Aaron von Kreisler", "email": "invalidemail.com", "password":"password"`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Register)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("Failed to return expected status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}
}

// TestRegisterFailWithBadPassword tests whether Register handler returns a
// status code of http.StatusUnprocessableEntity if the password doesn't meet
// domain requirements.
func TestRegisterFailWithBadPassword(t *testing.T) {
	var jsonStr = []byte(`"name":"Aaron von Kreisler", "email": "valid@email.com", "password":"short"`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Register)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("Failed to return expected status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}
}
