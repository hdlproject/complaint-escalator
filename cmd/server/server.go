package main

import (
	"complaint-escalator/internal/config"
	"complaint-escalator/internal/email"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// EmailRequest represents the JSON request structure for sending emails
type EmailRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// EmailResponse represents the JSON response structure
type EmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      string `json:"id,omitempty"`
}

// Server represents the HTTP server
type Server struct {
	config      *config.Config
	emailClient *email.EmailClient
	httpServer  *http.Server
}

// NewServer creates a new HTTP server instance
func NewServer(cfg config.Config) (*Server, error) {
	// Initialize email client
	emailClient, err := email.GetEmailClientFromConfig(cfg.ACS.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize email client: %w", err)
	}

	server := &Server{
		config:      &cfg,
		emailClient: emailClient,
	}

	// Create HTTP server
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/health", server.healthHandler)
	mux.HandleFunc("/email/send", server.sendEmailHandler)

	server.httpServer = &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("Starting HTTP server on port 8080")
	log.Printf("Available endpoints:")
	log.Printf("  GET  /health")
	log.Printf("  POST /email/send")

	return s.httpServer.ListenAndServe()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop() error {
	log.Println("Shutting down HTTP server...")
	return s.httpServer.Close()
}

// healthHandler handles health check requests
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "complaint-escalator-email-service",
	}

	json.NewEncoder(w).Encode(response)
}

// sendEmailHandler handles email sending requests
func (s *Server) sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var emailReq EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}
	if emailReq.Subject == "" {
		http.Error(w, "Subject is required", http.StatusBadRequest)
		return
	}
	if emailReq.Body == "" {
		http.Error(w, "Body is required", http.StatusBadRequest)
		return
	}

	emailMsg := email.CreateEmailMessageFromConfig(
		s.config.ACS.FromEmail,
		s.config.Email.To,
		s.config.Email.CC,
		s.config.Email.BCC,
		s.config.Email.ReplyTo,
		emailReq.Subject,
		emailReq.Body,
	)

	// Send email
	ctx := r.Context()
	err := s.emailClient.SendEmail(ctx, emailMsg)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		response := EmailResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to send email: %v", err),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	response := EmailResponse{
		Success: true,
		Message: "Email sent successfully",
		ID:      fmt.Sprintf("email_%d", time.Now().Unix()),
	}
	json.NewEncoder(w).Encode(response)
}

// middleware for logging and CORS
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)

		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
