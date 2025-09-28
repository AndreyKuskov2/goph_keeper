package config

// Flags represents command-line flags and environment variables for configuration.
type Flags struct {
	ConfigPath string `env:"CONFIG_PATH"` // Path to the configuration file
}

// Config represents the main configuration structure for the GophKeeper server.
type Config struct {
	Application *Application `yaml:"application"` // Application-specific configuration
	Database    *Database    `yaml:"database"`    // Database connection configuration
}

// Application represents application-specific configuration settings.
type Application struct {
	Address       string `yaml:"address"`          // Server address and port (e.g., ":8080")
	SecretToken   string `yaml:"jwt_secret_token"` // Secret key for JWT token signing
	Debug         bool   `yaml:"debug"`            // Enable debug mode for detailed logging
	AddSourcePath bool   `yaml:"add_source_path"`  // Include source file paths in log messages
}

// Database represents database connection configuration settings.
type Database struct {
	Type           string `yaml:"type"`            // Database type (e.g., "postgres")
	MigrationsPath string `yaml:"migrations_path"` // Path to database migration files
	Host           string `yaml:"host"`            // Database host address
	Port           int    `yaml:"port"`            // Database port number
	Name           string `yaml:"name"`            // Database name
	User           string `yaml:"user"`            // Database username
	Password       string `yaml:"password"`        // Database password
}
