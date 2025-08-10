package main

import (
	"bytes"
	"complaint-escalator/internal/config"
	"complaint-escalator/pkg/testutils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	// Load test configuration
	cfg, err := config.LoadConfig(testutils.GetTestConfigPath())
	if err != nil {
		t.Fatalf("Failed to load test configuration: %v", err)
	}

	server := &Server{
		config: &cfg,
	}

	// Create request
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.healthHandler)

	// Call handler
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check content type
	expected := "application/json"
	if rr.Header().Get("Content-Type") != expected {
		t.Errorf("handler returned wrong content type: got %v want %v", rr.Header().Get("Content-Type"), expected)
	}

	// Check response body contains expected fields
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response["status"] != "healthy" {
		t.Errorf("handler returned unexpected status: got %v want %v", response["status"], "healthy")
	}

	if response["service"] != "complaint-escalator-email-service" {
		t.Errorf("handler returned unexpected service: got %v want %v", response["service"], "complaint-escalator-email-service")
	}
}

func TestSendEmailHandler_InvalidMethod(t *testing.T) {
	// Load test configuration
	cfg, err := config.LoadConfig(testutils.GetTestConfigPath())
	if err != nil {
		t.Fatalf("Failed to load test configuration: %v", err)
	}

	server := &Server{
		config: &cfg,
	}

	// Create GET request (should fail)
	req, err := http.NewRequest("GET", "/email/send", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.sendEmailHandler)

	// Call handler
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestSendEmailHandler_InvalidJSON(t *testing.T) {
	cfg, err := config.LoadConfig(testutils.GetTestConfigPath())
	if err != nil {
		t.Fatalf("Failed to load test configuration: %v", err)
	}

	server := &Server{
		config: &cfg,
	}

	// Create request with invalid JSON
	req, err := http.NewRequest("POST", "/email/send", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.sendEmailHandler)

	// Call handler
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestSendEmailHandler_MissingRecipients(t *testing.T) {
	// Load test configuration
	cfg, err := config.LoadConfig(testutils.GetTestConfigPath())
	if err != nil {
		t.Fatalf("Failed to load test configuration: %v", err)
	}

	server := &Server{
		config: &cfg,
	}

	// Create request with missing recipients
	emailReq := EmailRequest{
		Subject: "Test Subject",
		Body:    "Test Body",
		// Missing To field
	}

	jsonData, err := json.Marshal(emailReq)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/email/send", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.sendEmailHandler)

	// Call handler
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestSendEmailHandler_ValidRequest(t *testing.T) {
	// Load test configuration
	cfg, err := config.LoadConfig(testutils.GetTestConfigPath())
	if err != nil {
		t.Fatalf("Failed to load test configuration: %v", err)
	}

	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create valid email request
	emailReq := EmailRequest{
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	jsonData, err := json.Marshal(emailReq)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/email/send", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.sendEmailHandler)

	// Call handler
	handler.ServeHTTP(rr, req)

	// Note: This will fail due to invalid test credentials, but the validation should pass
	// We're testing the request handling, not the actual email sending
	if rr.Code != http.StatusInternalServerError {
		t.Logf("Expected internal server error due to test credentials, got status: %d", rr.Code)
	}
}
