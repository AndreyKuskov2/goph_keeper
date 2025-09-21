package config

type Flags struct {
	ConfigPath string `env:"CONFIG_PATH"`
}

type Config struct {
	Application *Application `yaml:"application"`
	Database    *Database    `yaml:"database"`
}

type Application struct {
	Address       string `yaml:"address"`
	SecretToken   string `yaml:"jwt_secret_token"`
	Debug         bool   `yaml:"debug"`
	AddSourcePath bool   `yaml:"add_source_path"`
}

type Database struct {
	Type           string `yaml:"type"`
	MigrationsPath string `yaml:"migrations_path"`
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	Name           string `yaml:"name"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
}
