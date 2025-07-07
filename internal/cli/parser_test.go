package cli

import (
	"os/user"
	"strings"
	"testing"
)

func TestParseHost(t *testing.T) {
	tests := []struct {
		name         string
		arg          string
		expectedUser string
		expectedHost string
		expectError  bool
	}{
		{
			name:         "user@host format",
			arg:          "testuser@example.com",
			expectedUser: "testuser",
			expectedHost: "example.com",
			expectError:  false,
		},
		{
			name:         "host only",
			arg:          "example.com",
			expectedUser: "", // Will be current user
			expectedHost: "example.com",
			expectError:  false,
		},
		{
			name:        "multiple @ symbols",
			arg:         "user@host@extra",
			expectError: true,
		},
		{
			name:         "empty string",
			arg:          "",
			expectedUser: "", // Will be current user
			expectedHost: "",
			expectError:  false,
		},
		{
			name:         "IP address with user",
			arg:          "root@192.168.1.100",
			expectedUser: "root",
			expectedHost: "192.168.1.100",
			expectError:  false,
		},
		{
			name:         "IP address only",
			arg:          "192.168.1.100",
			expectedUser: "", // Will be current user
			expectedHost: "192.168.1.100",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username, host, err := ParseHost(tt.arg)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.expectedUser != "" {
				if username != tt.expectedUser {
					t.Errorf("Expected user %s, got %s", tt.expectedUser, username)
				}
			} else {
				// When no user is specified, should return current user
				currentUser, err := user.Current()
				if err != nil {
					t.Errorf("Failed to get current user: %v", err)
				} else if username != currentUser.Username {
					t.Errorf("Expected current user %s, got %s", currentUser.Username, username)
				}
			}

			if host != tt.expectedHost {
				t.Errorf("Expected host %s, got %s", tt.expectedHost, host)
			}
		})
	}
}

func TestParseHostErrorMessages(t *testing.T) {
	// Test specific error messages
	_, _, err := ParseHost("user@host@extra")
	if err == nil {
		t.Error("Expected error for multiple @ symbols")
	} else if !strings.Contains(err.Error(), "Illigal connection") {
		t.Errorf("Expected 'Illigal connection' error, got: %v", err)
	}
}