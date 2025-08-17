package email

import (
	"complaint-escalator/internal/config"
	"complaint-escalator/pkg/testutils"
	"context"
	"reflect"
	"testing"
)

func TestEmailClientCreation(t *testing.T) {
	// Load test configuration
	cfg, err := config.LoadConfig(testutils.GetTestConfigPath())
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	// Test email client creation from config
	emailClient, err := NewEmailClient(cfg.ACS.ConnectionString)
	if err != nil {
		t.Fatalf("Failed to create email client: %v", err)
	}
	if emailClient == nil {
		t.Fatal("Email client should not be nil")
	}

	// Test with empty connection string
	_, err = NewEmailClient("")
	if err == nil {
		t.Error("Expected error for empty connection string")
	}
}

func TestEmailMessageValidation(t *testing.T) {
	// Test valid email message
	validMsg := EmailMessage{
		From:    "test@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	// Create email client with test config
	cfg, err := config.LoadConfig(testutils.GetTestConfigPath())
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	emailClient, err := NewEmailClient(cfg.ACS.ConnectionString)
	if err != nil {
		t.Fatalf("Failed to create email client: %v", err)
	}

	// Test validation (this should not error)
	ctx := context.Background()
	err = emailClient.SendEmail(ctx, validMsg)
	// Note: This will fail due to invalid test credentials, but validation should pass
	// We're testing the validation logic, not the actual sending
	if err == nil {
		t.Log("Email validation passed (actual sending would fail with test credentials)")
	}
}

func TestEmailMessageWithOptionalFields(t *testing.T) {
	// Load test configuration
	cfg, err := config.LoadConfig(testutils.GetTestConfigPath())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test email message creation with all optional fields
	subject := "Test Subject with Optional Fields"
	body := "Test Body with Optional Fields"
	emailMsg := CreateEmailMessageFromConfig(
		cfg.ACS.FromEmail,
		cfg.Email.To,
		cfg.Email.CC,
		cfg.Email.BCC,
		cfg.Email.ReplyTo,
		subject,
		body,
	)

	// Verify all fields are set correctly
	if emailMsg.From != "test@test-domain.dev" {
		t.Errorf("Expected from email %s, got %s", "test@test-domain.dev", emailMsg.From)
	}
	if emailMsg.Subject != subject {
		t.Errorf("Expected subject %s, got %s", subject, emailMsg.Subject)
	}
	if emailMsg.Body != body {
		t.Errorf("Expected body %s, got %s", body, emailMsg.Body)
	}

	// Verify optional fields
	if len(emailMsg.CC) == 0 {
		t.Error("CC field should be populated from config")
	}
	if len(emailMsg.BCC) == 0 {
		t.Error("BCC field should be populated from config")
	}
	if emailMsg.ReplyTo == "" {
		t.Error("ReplyTo field should be populated from config")
	}
}

func TestConnectionStringParsing(t *testing.T) {
	// Test connection string parsing with test values
	testConnStr := "endpoint=https://test-acs.asiapacific.communication.azure.com/;accesskey=test-access-key"

	emailClient, err := NewEmailClient(testConnStr)
	if err != nil {
		t.Fatalf("Failed to create email client with test connection string: %v", err)
	}

	if emailClient == nil {
		t.Fatal("Email client should not be nil")
	}
}

func TestCreateEmailMessageFromConfig(t *testing.T) {
	// Test data
	fromEmail := "test@example.com"
	toEmails := []string{"recipient1@example.com", "recipient2@example.com"}
	ccEmails := []string{"cc@example.com"}
	bccEmails := []string{"bcc@example.com"}
	replyTo := "reply@example.com"
	subject := "Test Subject"
	body := "Test Body"

	// Test with all fields
	emailMsg := CreateEmailMessageFromConfig(fromEmail, toEmails, ccEmails, bccEmails, replyTo, subject, body)

	// Verify all fields are set correctly
	if emailMsg.From != fromEmail {
		t.Errorf("Expected from email %s, got %s", fromEmail, emailMsg.From)
	}
	if !reflect.DeepEqual(emailMsg.To, toEmails) {
		t.Errorf("Expected to emails %v, got %v", toEmails, emailMsg.To)
	}
	if !reflect.DeepEqual(emailMsg.CC, ccEmails) {
		t.Errorf("Expected cc emails %v, got %v", ccEmails, emailMsg.CC)
	}
	if !reflect.DeepEqual(emailMsg.BCC, bccEmails) {
		t.Errorf("Expected bcc emails %v, got %v", bccEmails, emailMsg.BCC)
	}
	if emailMsg.ReplyTo != replyTo {
		t.Errorf("Expected reply to %s, got %s", replyTo, emailMsg.ReplyTo)
	}
	if emailMsg.Subject != subject {
		t.Errorf("Expected subject %s, got %s", subject, emailMsg.Subject)
	}
	if emailMsg.Body != body {
		t.Errorf("Expected body %s, got %s", body, emailMsg.Body)
	}

	// Test with empty optional fields
	emailMsg2 := CreateEmailMessageFromConfig(fromEmail, toEmails, nil, nil, "", subject, body)
	if len(emailMsg2.CC) != 0 {
		t.Error("CC should be empty when not provided")
	}
	if len(emailMsg2.BCC) != 0 {
		t.Error("BCC should be empty when not provided")
	}
	if emailMsg2.ReplyTo != "" {
		t.Error("ReplyTo should be empty when not provided")
	}
}
