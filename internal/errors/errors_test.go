package errors

import (
	"errors"
	"fmt"
	"testing"
)

// TestConfigError tests the ConfigError type
func TestConfigError(t *testing.T) {
	tests := []struct {
		name          string
		filePath      string
		operation     string
		field         string
		err           error
		constructor   func(error) error
		expectedMsg   string
		shouldContain string
	}{
		{
			name:          "ConfigError with operation",
			filePath:      "/etc/config.json",
			operation:     "parsing",
			err:           fmt.Errorf("invalid JSON"),
			constructor:   func(err error) error { return NewConfigError("/etc/config.json", "parsing", err) },
			expectedMsg:   "config error in /etc/config.json: parsing - invalid JSON",
			shouldContain: "config error",
		},
		{
			name:          "ConfigError with field",
			filePath:      "config.json",
			field:         "msm_path",
			err:           fmt.Errorf("required field missing"),
			constructor:   func(err error) error { return NewConfigFieldError("config.json", "msm_path", err) },
			expectedMsg:   "config error in config.json: field 'msm_path' - required field missing",
			shouldContain: "field 'msm_path'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constructor(tt.err)
			msg := err.Error()

			if msg != tt.expectedMsg {
				t.Errorf("Expected message:\n%s\nGot:\n%s", tt.expectedMsg, msg)
			}

			if !contains(msg, tt.shouldContain) {
				t.Errorf("Message should contain '%s', got: %s", tt.shouldContain, msg)
			}

			// Test Cause method
			configErr, ok := err.(*ConfigError)
			if !ok {
				t.Errorf("Expected *ConfigError, got %T", err)
			}
			if configErr.Cause() == nil {
				t.Errorf("Cause() should not be nil")
			}
		})
	}
}

// TestPropertiesError tests the PropertiesError type
func TestPropertiesError(t *testing.T) {
	tests := []struct {
		name          string
		filePath      string
		operation     string
		err           error
		shouldContain string
	}{
		{
			name:          "PropertiesError on read",
			filePath:      "/opt/minecraft/server.properties",
			operation:     "read",
			err:           fmt.Errorf("permission denied"),
			shouldContain: "properties file error",
		},
		{
			name:          "PropertiesError on write",
			filePath:      "server.properties",
			operation:     "write",
			err:           fmt.Errorf("disk full"),
			shouldContain: "write",
		},
		{
			name:          "PropertiesError on parse",
			filePath:      "server.properties",
			operation:     "parse",
			err:           fmt.Errorf("invalid format"),
			shouldContain: "parse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewPropertiesError(tt.filePath, tt.operation, tt.err)
			msg := err.Error()

			if !contains(msg, tt.shouldContain) {
				t.Errorf("Message should contain '%s', got: %s", tt.shouldContain, msg)
			}

			if err.Cause() != tt.err {
				t.Errorf("Cause() should return the original error")
			}
		})
	}
}

// TestWorldError tests the WorldError type
func TestWorldError(t *testing.T) {
	tests := []struct {
		name          string
		worldName     string
		operation     string
		err           error
		shouldContain string
	}{
		{
			name:          "WorldError on create",
			worldName:     "survival",
			operation:     "create",
			err:           fmt.Errorf("MSM not found"),
			shouldContain: "world error (create 'survival')",
		},
		{
			name:          "WorldError on delete",
			worldName:     "creative",
			operation:     "delete",
			err:           fmt.Errorf("world directory not found"),
			shouldContain: "delete",
		},
		{
			name:          "WorldError on restore",
			worldName:     "nether",
			operation:     "restore",
			err:           fmt.Errorf("backup not found"),
			shouldContain: "restore",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewWorldError(tt.worldName, tt.operation, tt.err)
			msg := err.Error()

			if !contains(msg, tt.shouldContain) {
				t.Errorf("Message should contain '%s', got: %s", tt.shouldContain, msg)
			}

			if err.Cause() != tt.err {
				t.Errorf("Cause() should return the original error")
			}
		})
	}
}

// TestBuildError tests the BuildError type
func TestBuildError(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		phase         string
		err           error
		shouldContain string
	}{
		{
			name:          "BuildError with phase",
			version:       "1.20.1",
			phase:         "compilation",
			err:           fmt.Errorf("Java compilation failed"),
			shouldContain: "build error (version 1.20.1, phase 'compilation')",
		},
		{
			name:          "BuildError without phase",
			version:       "1.19",
			phase:         "",
			err:           fmt.Errorf("BuildTools not found"),
			shouldContain: "build error (version 1.19",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewBuildError(tt.version, tt.phase, tt.err)
			msg := err.Error()

			if !contains(msg, tt.shouldContain) {
				t.Errorf("Message should contain '%s', got: %s", tt.shouldContain, msg)
			}

			if err.Cause() != tt.err {
				t.Errorf("Cause() should return the original error")
			}
		})
	}
}

// TestPublishError tests the PublishError type
func TestPublishError(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		phase         string
		err           error
		shouldContain string
	}{
		{
			name:          "PublishError with phase",
			filename:      "spigot-1.20.1.jar",
			phase:         "git push",
			err:           fmt.Errorf("authentication failed"),
			shouldContain: "publish error (phase 'git push'",
		},
		{
			name:          "PublishError without phase",
			filename:      "spigot-1.19.jar",
			phase:         "",
			err:           fmt.Errorf("network timeout"),
			shouldContain: "spigot-1.19.jar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewPublishError(tt.filename, tt.phase, tt.err)
			msg := err.Error()

			if !contains(msg, tt.shouldContain) {
				t.Errorf("Message should contain '%s', got: %s", tt.shouldContain, msg)
			}

			if err.Cause() != tt.err {
				t.Errorf("Cause() should return the original error")
			}
		})
	}
}

// TestSubprocessError tests the SubprocessError type
func TestSubprocessError(t *testing.T) {
	tests := []struct {
		name          string
		command       string
		reason        string
		exitCode      int
		stderr        string
		shouldContain string
	}{
		{
			name:          "SubprocessError with exit code",
			command:       "msm jar create",
			reason:        "",
			exitCode:      1,
			stderr:        "jar file not found",
			shouldContain: "subprocess error executing 'msm jar create': exit code 1",
		},
		{
			name:          "SubprocessError with reason",
			command:       "java",
			reason:        "command not found",
			exitCode:      0,
			stderr:        "",
			shouldContain: "command not found",
		},
		{
			name:          "SubprocessError with exit code no stderr",
			command:       "git push",
			reason:        "",
			exitCode:      128,
			stderr:        "",
			shouldContain: "exit code 128",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.exitCode != 0 {
				err = NewSubprocessExitError(tt.command, tt.exitCode, tt.stderr)
			} else {
				err = NewSubprocessError(tt.command, tt.reason, fmt.Errorf("test error"))
			}

			msg := err.Error()
			if !contains(msg, tt.shouldContain) {
				t.Errorf("Message should contain '%s', got: %s", tt.shouldContain, msg)
			}

			subErr := err.(*SubprocessError)
			if subErr.Cause() == nil {
				t.Errorf("Cause() should not be nil")
			}
		})
	}
}

// TestFileSystemError tests the FileSystemError type
func TestFileSystemError(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		operation     string
		err           error
		shouldContain string
	}{
		{
			name:          "FileSystemError on read",
			path:          "/opt/minecraft/server.jar",
			operation:     "read",
			err:           fmt.Errorf("no such file"),
			shouldContain: "filesystem error (read)",
		},
		{
			name:          "FileSystemError on write",
			path:          "/opt/config",
			operation:     "write",
			err:           fmt.Errorf("permission denied"),
			shouldContain: "write",
		},
		{
			name:          "FileSystemError on delete",
			path:          "/opt/minecraft/old/",
			operation:     "delete",
			err:           fmt.Errorf("directory not empty"),
			shouldContain: "delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewFileSystemError(tt.path, tt.operation, tt.err)
			msg := err.Error()

			if !contains(msg, tt.shouldContain) {
				t.Errorf("Message should contain '%s', got: %s", tt.shouldContain, msg)
			}

			if err.Cause() != tt.err {
				t.Errorf("Cause() should return the original error")
			}
		})
	}
}

// TestErrorsImplementError verifies all error types implement the error interface
func TestErrorsImplementError(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "ConfigError implements error",
			err:  NewConfigError("test.json", "test", fmt.Errorf("test")),
		},
		{
			name: "PropertiesError implements error",
			err:  NewPropertiesError("props.txt", "read", fmt.Errorf("test")),
		},
		{
			name: "WorldError implements error",
			err:  NewWorldError("world", "create", fmt.Errorf("test")),
		},
		{
			name: "BuildError implements error",
			err:  NewBuildError("1.20", "compile", fmt.Errorf("test")),
		},
		{
			name: "PublishError implements error",
			err:  NewPublishError("jar.jar", "push", fmt.Errorf("test")),
		},
		{
			name: "SubprocessError implements error",
			err:  NewSubprocessError("cmd", "reason", fmt.Errorf("test")),
		},
		{
			name: "FileSystemError implements error",
			err:  NewFileSystemError("/path", "read", fmt.Errorf("test")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This will fail to compile if the type doesn't implement error
			var _ error = tt.err
		})
	}
}

// TestErrorWrapping tests that errors can be wrapped and retrieved
func TestErrorWrapping(t *testing.T) {
	originalErr := fmt.Errorf("original error")

	tests := []struct {
		name         string
		wrappedErr   interface{ Cause() error }
		expectedCause error
	}{
		{
			name:         "ConfigError wrapping",
			wrappedErr:   NewConfigError("file", "op", originalErr),
			expectedCause: originalErr,
		},
		{
			name:         "WorldError wrapping",
			wrappedErr:   NewWorldError("name", "op", originalErr),
			expectedCause: originalErr,
		},
		{
			name:         "BuildError wrapping",
			wrappedErr:   NewBuildError("ver", "phase", originalErr),
			expectedCause: originalErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cause := tt.wrappedErr.Cause()
			if cause != tt.expectedCause {
				t.Errorf("Expected cause to be the original error, got %v", cause)
			}
		})
	}
}

// TestSubprocessExitErrorFormat tests exit code error formatting
func TestSubprocessExitErrorFormat(t *testing.T) {
	tests := []struct {
		name          string
		command       string
		exitCode      int
		stderr        string
		expectedMsg   string
	}{
		{
			name:          "Exit code 1 with stderr",
			command:       "test-cmd",
			exitCode:      1,
			stderr:        "Error: file not found",
			expectedMsg:   "subprocess error executing 'test-cmd': exit code 1 - Error: file not found",
		},
		{
			name:          "Exit code 127 without stderr",
			command:       "missing-cmd",
			exitCode:      127,
			stderr:        "",
			expectedMsg:   "subprocess error executing 'missing-cmd': exit code 127",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewSubprocessExitError(tt.command, tt.exitCode, tt.stderr)
			msg := err.Error()

			if msg != tt.expectedMsg {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expectedMsg, msg)
			}
		})
	}
}

// TestErrorChaining tests that errors properly chain using Go's error interface
func TestErrorChaining(t *testing.T) {
	originalErr := fmt.Errorf("original problem")
	configErr := NewConfigError("config.json", "loading", originalErr)

	// Test that we can use errors.Is with the cause
	if !errors.Is(configErr.Cause(), originalErr) {
		t.Errorf("Error chaining failed: cause should match original error")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
