package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response.
	rr := httptest.NewRecorder()

	// Llama al handler
	handler := http.HandlerFunc(HealthCheckHandler)
	handler.ServeHTTP(rr, req)

	// Verify status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler devolvió código de estado incorrecto: obtuvo %v, se esperaba %v", status, http.StatusOK)
	}

	// Check that the body of the answer is ‘OK’.
	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler devolvió un cuerpo inesperado: obtuvo %v, se esperaba %v", rr.Body.String(), expected)
	}
}
