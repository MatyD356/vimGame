package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleHealt(t *testing.T) {
	// Create a request to pass to the handler (no body needed)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler with the ResponseRecorder and Request
	HandleHealt(rr, req)

	// Check the status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect
	expected := "Server is healthy\n" // fmt.Fprintln adds newline
	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(expected) {
		t.Errorf("Handler returned unexpected body: got %q want %q", rr.Body.String(), expected)
	}
}
