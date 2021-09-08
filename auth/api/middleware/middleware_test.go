package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP handler that we can pass to our secureHeaders middleware
	// which writes a 200 status code and "OK" response
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Pass the mock HTTP handler to the middleware. Because SecureHEaders *returns*
	// a http.Handler we can call its ServerHTTP() method, passing in the
	// http.ResponseRecorder and dummy http.Request to execute it

	SecureHeaders(next).ServeHTTP(rr, r)

	// Call the Result() method to get the results
	rs := rr.Result()

	// Check that the middleware has correctly set the X-Frame-Options
	// on the response
	frameOptions := rs.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}

	// Check that the middleware has correctly set the X-XSS-Protection header
	// on the response
	xssProtection := rs.Header.Get("X-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf("want %q; got %q", "1; mode=block", xssProtection)
	}
	// Check that the middleware has correctly called the next handler in line
	// and the response status code and body are as expected.
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to ewuad %q", "OK")
	}

}
