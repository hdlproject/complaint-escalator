package config

import (
	"complaint-escalator/pkg/testutils"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Load test configuration
	cfg, err := LoadConfig(testutils.GetTestConfigPath())
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	// Test interval configuration
	expectedInterval := 5 * time.Minute
	if cfg.Interval != expectedInterval {
		t.Errorf("Expected interval %v, got %v", expectedInterval, cfg.Interval)
	}

	// Test backoff configuration
	expectedBackoff := 2 * time.Minute
	if cfg.Backoff != expectedBackoff {
		t.Errorf("Expected backoff %v, got %v", expectedBackoff, cfg.Backoff)
	}

	// Test template configuration
	expectedTemplate := "Test complaint template for automated testing."
	if cfg.Template != expectedTemplate {
		t.Errorf("Expected template %s, got %s", expectedTemplate, cfg.Template)
	}

	// Test channels configuration
	expectedChannels := []string{"email", "notification"}
	if len(cfg.Channels) != len(expectedChannels) {
		t.Errorf("Expected %d channels, got %d", len(expectedChannels), len(cfg.Channels))
	}
	for i, channel := range expectedChannels {
		if cfg.Channels[i] != channel {
			t.Errorf("Expected channel %s at position %d, got %s", channel, i, cfg.Channels[i])
		}
	}

	// Test ACS configuration
	if cfg.ACS.ConnectionString != "endpoint=https://test-acs.asiapacific.communication.azure.com/;accesskey=test-access-key" {
		t.Error("ACS connection string not loaded correctly")
	}
	if cfg.ACS.Domain != "test-domain.dev" {
		t.Error("ACS domain not loaded correctly")
	}
	if cfg.ACS.FromEmail != "test@test-domain.dev" {
		t.Error("ACS from email not loaded correctly")
	}

	// Test email configuration
	if len(cfg.Email.To) != 1 || cfg.Email.To[0] != "test@example.com" {
		t.Error("Email 'to' not loaded correctly")
	}
	if len(cfg.Email.CC) != 1 || cfg.Email.CC[0] != "cc-test@example.com" {
		t.Error("Email 'cc' not loaded correctly")
	}
	if len(cfg.Email.BCC) != 1 || cfg.Email.BCC[0] != "bcc-test@example.com" {
		t.Error("Email 'bcc' not loaded correctly")
	}
	if cfg.Email.ReplyTo != "reply-test@example.com" {
		t.Error("Email 'reply_to' not loaded correctly")
	}
}
