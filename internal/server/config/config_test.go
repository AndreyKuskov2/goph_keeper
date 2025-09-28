package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name          string
		configPath    string
		configContent string
		wantErr       bool
		setupFile     bool
	}{
		{
			name:       "valid config file",
			configPath: "test_config.yml",
			configContent: `
application:
  address: ":8080"
  jwt_secret_token: "test-secret"
  debug: true
  add_source_path: false
database:
  type: "postgres"
  migrations_path: "./migrations"
  host: "localhost"
  port: 5432
  name: "test_db"
  user: "test_user"
  password: "test_password"
`,
			wantErr:   false,
			setupFile: true,
		},
		{
			name:       "invalid yaml",
			configPath: "invalid_config.yml",
			configContent: `
application:
  address: ":8080"
  jwt_secret_token: "test-secret"
  debug: true
  add_source_path: false
database:
  type: "postgres"
  migrations_path: "./migrations"
  host: "localhost"
  port: 5432
  name: "test_db"
  user: "test_user"
  password: "test_password"
invalid: yaml: content
`,
			wantErr:   true,
			setupFile: true,
		},
		{
			name:       "file not found",
			configPath: "nonexistent_config.yml",
			wantErr:    true,
			setupFile:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test file if needed
			if tt.setupFile {
				tempDir := t.TempDir()
				configFile := filepath.Join(tempDir, tt.configPath)
				err := os.WriteFile(configFile, []byte(tt.configContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create test config file: %v", err)
				}

				// Set environment variable for config path
				os.Setenv("CONFIG_PATH", configFile)
				defer os.Unsetenv("CONFIG_PATH")
			}

			// Reset command line flags
			os.Args = []string{"test"}
			// Clear any existing flags
			pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

			cfg, err := NewConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && cfg == nil {
				t.Error("Expected config to be created")
			}

			if !tt.wantErr && cfg != nil {
				// Verify config structure
				if cfg.Application == nil {
					t.Error("Expected Application config to be set")
				}
				if cfg.Database == nil {
					t.Error("Expected Database config to be set")
				}
			}
		})
	}
}

func TestConfigStruct(t *testing.T) {
	cfg := &Config{
		Application: &Application{
			Address:       ":8080",
			SecretToken:   "test-secret",
			Debug:         true,
			AddSourcePath: false,
		},
		Database: &Database{
			Type:           "postgres",
			MigrationsPath: "./migrations",
			Host:           "localhost",
			Port:           5432,
			Name:           "test_db",
			User:           "test_user",
			Password:       "test_password",
		},
	}

	// Test Application fields
	if cfg.Application.Address != ":8080" {
		t.Errorf("Expected Address :8080, got %s", cfg.Application.Address)
	}
	if cfg.Application.SecretToken != "test-secret" {
		t.Errorf("Expected SecretToken test-secret, got %s", cfg.Application.SecretToken)
	}
	if !cfg.Application.Debug {
		t.Error("Expected Debug to be true")
	}
	if cfg.Application.AddSourcePath {
		t.Error("Expected AddSourcePath to be false")
	}

	// Test Database fields
	if cfg.Database.Type != "postgres" {
		t.Errorf("Expected Type postgres, got %s", cfg.Database.Type)
	}
	if cfg.Database.MigrationsPath != "./migrations" {
		t.Errorf("Expected MigrationsPath ./migrations, got %s", cfg.Database.MigrationsPath)
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected Host localhost, got %s", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("Expected Port 5432, got %d", cfg.Database.Port)
	}
	if cfg.Database.Name != "test_db" {
		t.Errorf("Expected Name test_db, got %s", cfg.Database.Name)
	}
	if cfg.Database.User != "test_user" {
		t.Errorf("Expected User test_user, got %s", cfg.Database.User)
	}
	if cfg.Database.Password != "test_password" {
		t.Errorf("Expected Password test_password, got %s", cfg.Database.Password)
	}
}

func TestFlagsStruct(t *testing.T) {
	flags := &Flags{
		ConfigPath: "/path/to/config.yml",
	}

	if flags.ConfigPath != "/path/to/config.yml" {
		t.Errorf("Expected ConfigPath /path/to/config.yml, got %s", flags.ConfigPath)
	}
}

func TestNewConfig_EnvironmentVariable(t *testing.T) {
	// Create test config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "env_test_config.yml")
	configContent := `
application:
  address: ":8080"
  jwt_secret_token: "test-secret"
  debug: true
  add_source_path: false
database:
  type: "postgres"
  migrations_path: "./migrations"
  host: "localhost"
  port: 5432
  name: "test_db"
  user: "test_user"
  password: "test_password"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Set environment variable
	os.Setenv("CONFIG_PATH", configFile)
	defer os.Unsetenv("CONFIG_PATH")

	// Reset command line flags
	os.Args = []string{"test"}

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}

	if cfg == nil {
		t.Error("Expected config to be created")
	}

	if cfg.Application == nil {
		t.Error("Expected Application config to be set")
	}

	if cfg.Database == nil {
		t.Error("Expected Database config to be set")
	}
}
