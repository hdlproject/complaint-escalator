package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// EmailClient represents an Azure Communication Services email client
type EmailClient struct {
	connectionString string
	httpClient       *http.Client
	endpoint         string
	accessKey        string
}

// EmailMessage represents an email message to be sent
type EmailMessage struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	// Optional fields
	CC      []string `json:"cc,omitempty"`
	BCC     []string `json:"bcc,omitempty"`
	ReplyTo string   `json:"replyTo,omitempty"`
}

// Azure Email API Request structure
type azureEmailRequest struct {
	SenderAddress string `json:"senderAddress"`
	Content       struct {
		Subject   string `json:"subject"`
		PlainText string `json:"plainText"`
	} `json:"content"`
	Recipients struct {
		To  []azureEmailAddress `json:"to"`
		CC  []azureEmailAddress `json:"cc,omitempty"`
		BCC []azureEmailAddress `json:"bcc,omitempty"`
	} `json:"recipients"`
	ReplyTo *azureEmailAddress `json:"replyTo,omitempty"`
}

type azureEmailAddress struct {
	Email string `json:"email"`
}

// NewEmailClient creates a new email client using the provided connection string
func NewEmailClient(connectionString string) (*EmailClient, error) {
	if connectionString == "" {
		return nil, fmt.Errorf("connection string cannot be empty")
	}

	// Parse connection string
	endpoint, accessKey, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	return &EmailClient{
		connectionString: connectionString,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		endpoint:  endpoint,
		accessKey: accessKey,
	}, nil
}

// SendEmail sends an email using the Azure Communication Services REST API
func (ec *EmailClient) SendEmail(ctx context.Context, msg EmailMessage) error {
	if err := ec.validateMessage(msg); err != nil {
		return fmt.Errorf("invalid email message: %w", err)
	}

	// Convert to Azure API format
	azureReq := azureEmailRequest{
		SenderAddress: msg.From,
		Content: struct {
			Subject   string `json:"subject"`
			PlainText string `json:"plainText"`
		}{
			Subject:   msg.Subject,
			PlainText: msg.Body,
		},
		Recipients: struct {
			To  []azureEmailAddress `json:"to"`
			CC  []azureEmailAddress `json:"cc,omitempty"`
			BCC []azureEmailAddress `json:"bcc,omitempty"`
		}{
			To: convertToAzureEmailAddresses(msg.To),
		},
	}

	// Add CC recipients if provided
	if len(msg.CC) > 0 {
		azureReq.Recipients.CC = convertToAzureEmailAddresses(msg.CC)
	}

	// Add BCC recipients if provided
	if len(msg.BCC) > 0 {
		azureReq.Recipients.BCC = convertToAzureEmailAddresses(msg.BCC)
	}

	// Add Reply-To if provided
	if msg.ReplyTo != "" {
		azureReq.ReplyTo = &azureEmailAddress{Email: msg.ReplyTo}
	}

	// Convert to JSON
	jsonData, err := json.Marshal(azureReq)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/emails:send?api-version=2023-03-31", ec.endpoint)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("HMAC-SHA256 %s", ec.accessKey))

	// Send request
	resp, err := ec.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status
	if resp.StatusCode >= 400 {
		return fmt.Errorf("email send failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Email sent successfully. Status: %d", resp.StatusCode)
	return nil
}

// validateMessage validates the email message
func (ec *EmailClient) validateMessage(msg EmailMessage) error {
	if msg.From == "" {
		return fmt.Errorf("sender address is required")
	}
	if len(msg.To) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}
	if msg.Subject == "" {
		return fmt.Errorf("subject is required")
	}
	if msg.Body == "" {
		return fmt.Errorf("body is required")
	}
	return nil
}

// parseConnectionString parses the Azure Communication Services connection string
func parseConnectionString(connStr string) (endpoint, accessKey string, err error) {
	parts := strings.Split(connStr, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, "endpoint=") {
			endpoint = strings.TrimPrefix(part, "endpoint=")
		} else if strings.HasPrefix(part, "accesskey=") {
			accessKey = strings.TrimPrefix(part, "accesskey=")
		}
	}

	if endpoint == "" || accessKey == "" {
		return "", "", fmt.Errorf("invalid connection string format")
	}

	return endpoint, accessKey, nil
}

// convertToAzureEmailAddresses converts string slices to Azure email address slices
func convertToAzureEmailAddresses(addresses []string) []azureEmailAddress {
	result := make([]azureEmailAddress, len(addresses))
	for i, addr := range addresses {
		result[i] = azureEmailAddress{Email: addr}
	}
	return result
}

// GetEmailClientFromConfig creates an email client from the provided configuration
func GetEmailClientFromConfig(acsConnectionString string) (*EmailClient, error) {
	if acsConnectionString == "" {
		return nil, fmt.Errorf("ACS connection string is not configured")
	}

	client, err := NewEmailClient(acsConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create email client: %w", err)
	}

	return client, nil
}

// CreateEmailMessageFromConfig creates an email message from the provided configuration
func CreateEmailMessageFromConfig(fromEmail string, toEmails, ccEmails, bccEmails []string, replyTo string, subject, body string) EmailMessage {
	msg := EmailMessage{
		From:    fromEmail,
		To:      toEmails,
		Subject: subject,
		Body:    body,
		ReplyTo: replyTo,
	}

	// Add CC if configured
	if len(ccEmails) > 0 {
		msg.CC = ccEmails
	}

	// Add BCC if configured
	if len(bccEmails) > 0 {
		msg.BCC = bccEmails
	}

	return msg
}
