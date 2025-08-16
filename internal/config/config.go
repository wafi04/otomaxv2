package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server Configuration
	Server ServerConfig `mapstructure:"server"`
	
	// Database Configuration
	Database DatabaseConfig `mapstructure:"database"`
	
	// Redis Configuration
	Redis RedisConfig `mapstructure:"redis"`
	
	// JWT Configuration
	JWT JWTConfig `mapstructure:"jwt"`
	
	// Payment Gateway Configuration
	PaymentGateway PaymentGatewayConfig `mapstructure:"payment_gateway"`
	
	// External API Configuration
	ExternalAPI ExternalAPIConfig `mapstructure:"external_api"`
	
	// Logging Configuration
	Logging LoggingConfig `mapstructure:"logging"`
	
	// Application Configuration
	App AppConfig `mapstructure:"app"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"` // debug, release
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"db_name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	Timezone        string        `mapstructure:"timezone"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	MigrationPath   string        `mapstructure:"migration_path"`
}

type RedisConfig struct {
	Host        string        `mapstructure:"host"`
	Port        string        `mapstructure:"port"`
	Password    string        `mapstructure:"password"`
	DB          int           `mapstructure:"db"`
	MaxRetries  int           `mapstructure:"max_retries"`
	PoolSize    int           `mapstructure:"pool_size"`
	PoolTimeout time.Duration `mapstructure:"pool_timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

type JWTConfig struct {
	SecretKey      string        `mapstructure:"secret_key"`
	ExpireDuration time.Duration `mapstructure:"expire_duration"`
	RefreshDuration time.Duration `mapstructure:"refresh_duration"`
	Issuer         string        `mapstructure:"issuer"`
}

type PaymentGatewayConfig struct {
	Midtrans MidtransConfig `mapstructure:"midtrans"`
	Xendit   XenditConfig   `mapstructure:"xendit"`
	GoPay    GoPayConfig    `mapstructure:"gopay"`
}

type MidtransConfig struct {
	ServerKey    string `mapstructure:"server_key"`
	ClientKey    string `mapstructure:"client_key"`
	MerchantID   string `mapstructure:"merchant_id"`
	Environment  string `mapstructure:"environment"` // sandbox, production
	NotificationURL string `mapstructure:"notification_url"`
	ReturnURL    string `mapstructure:"return_url"`
	UnfinishURL  string `mapstructure:"unfinish_url"`
	ErrorURL     string `mapstructure:"error_url"`
}

type XenditConfig struct {
	SecretKey   string `mapstructure:"secret_key"`
	CallbackToken string `mapstructure:"callback_token"`
	Environment string `mapstructure:"environment"`
	WebhookURL  string `mapstructure:"webhook_url"`
}

type GoPayConfig struct {
	MerchantID  string `mapstructure:"merchant_id"`
	SecretKey   string `mapstructure:"secret_key"`
	Environment string `mapstructure:"environment"`
	CallbackURL string `mapstructure:"callback_url"`
}

type ExternalAPIConfig struct {
	Telkomsel TelkomselConfig `mapstructure:"telkomsel"`
	Indosat   IndosatConfig   `mapstructure:"indosat"`
	XL        XLConfig        `mapstructure:"xl"`
	Steam     SteamConfig     `mapstructure:"steam"`
	Garena    GarenaConfig    `mapstructure:"garena"`
}

type TelkomselConfig struct {
	BaseURL   string `mapstructure:"base_url"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	APIKey    string `mapstructure:"api_key"`
	Timeout   time.Duration `mapstructure:"timeout"`
}

type IndosatConfig struct {
	BaseURL   string `mapstructure:"base_url"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	APIKey    string `mapstructure:"api_key"`
	Timeout   time.Duration `mapstructure:"timeout"`
}

type XLConfig struct {
	BaseURL   string `mapstructure:"base_url"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	APIKey    string `mapstructure:"api_key"`
	Timeout   time.Duration `mapstructure:"timeout"`
}

type SteamConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type GarenaConfig struct {
	BaseURL   string `mapstructure:"base_url"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	APIKey    string `mapstructure:"api_key"`
	Timeout   time.Duration `mapstructure:"timeout"`
}

type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"` // json, text
	Output   string `mapstructure:"output"` // stdout, file
	FilePath string `mapstructure:"file_path"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"` // development, staging, production
	Debug       bool   `mapstructure:"debug"`
	URL         string `mapstructure:"url"`
	SecretKey   string `mapstructure:"secret_key"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "localhost"),
			Port:         getEnv("SERVER_PORT", "8080"),
			Mode:         getEnv("SERVER_MODE", "debug"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			Username:        getEnv("DB_USERNAME", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			DBName:          getEnv("DB_NAME", "topup_db"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			Timezone:        getEnv("DB_TIMEZONE", "Asia/Jakarta"),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 1*time.Hour),
			ConnMaxIdleTime: getDurationEnv("DB_CONN_MAX_IDLE_TIME", 10*time.Minute),
			MigrationPath:   getEnv("DB_MIGRATION_PATH", "internal/infrastructure/database/postgres/migrations"),
		},
		Redis: RedisConfig{
			Host:        getEnv("REDIS_HOST", "localhost"),
			Port:        getEnv("REDIS_PORT", "6379"),
			Password:    getEnv("REDIS_PASSWORD", ""),
			DB:          getIntEnv("REDIS_DB", 0),
			MaxRetries:  getIntEnv("REDIS_MAX_RETRIES", 3),
			PoolSize:    getIntEnv("REDIS_POOL_SIZE", 10),
			PoolTimeout: getDurationEnv("REDIS_POOL_TIMEOUT", 30*time.Second),
			IdleTimeout: getDurationEnv("REDIS_IDLE_TIMEOUT", 5*time.Minute),
		},
		JWT: JWTConfig{
			SecretKey:       getEnv("JWT_SECRET_KEY", "your-secret-key"),
			ExpireDuration:  getDurationEnv("JWT_EXPIRE_DURATION", 24*time.Hour),
			RefreshDuration: getDurationEnv("JWT_REFRESH_DURATION", 7*24*time.Hour),
			Issuer:          getEnv("JWT_ISSUER", "topup-app"),
		},
		PaymentGateway: PaymentGatewayConfig{
			Midtrans: MidtransConfig{
				ServerKey:       getEnv("MIDTRANS_SERVER_KEY", ""),
				ClientKey:       getEnv("MIDTRANS_CLIENT_KEY", ""),
				MerchantID:      getEnv("MIDTRANS_MERCHANT_ID", ""),
				Environment:     getEnv("MIDTRANS_ENVIRONMENT", "sandbox"),
				NotificationURL: getEnv("MIDTRANS_NOTIFICATION_URL", ""),
				ReturnURL:       getEnv("MIDTRANS_RETURN_URL", ""),
				UnfinishURL:     getEnv("MIDTRANS_UNFINISH_URL", ""),
				ErrorURL:        getEnv("MIDTRANS_ERROR_URL", ""),
			},
			Xendit: XenditConfig{
				SecretKey:     getEnv("XENDIT_SECRET_KEY", ""),
				CallbackToken: getEnv("XENDIT_CALLBACK_TOKEN", ""),
				Environment:   getEnv("XENDIT_ENVIRONMENT", "test"),
				WebhookURL:    getEnv("XENDIT_WEBHOOK_URL", ""),
			},
			GoPay: GoPayConfig{
				MerchantID:  getEnv("GOPAY_MERCHANT_ID", ""),
				SecretKey:   getEnv("GOPAY_SECRET_KEY", ""),
				Environment: getEnv("GOPAY_ENVIRONMENT", "development"),
				CallbackURL: getEnv("GOPAY_CALLBACK_URL", ""),
			},
		},
		ExternalAPI: ExternalAPIConfig{
			Telkomsel: TelkomselConfig{
				BaseURL:  getEnv("TELKOMSEL_BASE_URL", ""),
				Username: getEnv("TELKOMSEL_USERNAME", ""),
				Password: getEnv("TELKOMSEL_PASSWORD", ""),
				APIKey:   getEnv("TELKOMSEL_API_KEY", ""),
				Timeout:  getDurationEnv("TELKOMSEL_TIMEOUT", 30*time.Second),
			},
			Indosat: IndosatConfig{
				BaseURL:  getEnv("INDOSAT_BASE_URL", ""),
				Username: getEnv("INDOSAT_USERNAME", ""),
				Password: getEnv("INDOSAT_PASSWORD", ""),
				APIKey:   getEnv("INDOSAT_API_KEY", ""),
				Timeout:  getDurationEnv("INDOSAT_TIMEOUT", 30*time.Second),
			},
			XL: XLConfig{
				BaseURL:  getEnv("XL_BASE_URL", ""),
				Username: getEnv("XL_USERNAME", ""),
				Password: getEnv("XL_PASSWORD", ""),
				APIKey:   getEnv("XL_API_KEY", ""),
				Timeout:  getDurationEnv("XL_TIMEOUT", 30*time.Second),
			},
			Steam: SteamConfig{
				BaseURL: getEnv("STEAM_BASE_URL", ""),
				APIKey:  getEnv("STEAM_API_KEY", ""),
				Timeout: getDurationEnv("STEAM_TIMEOUT", 30*time.Second),
			},
			Garena: GarenaConfig{
				BaseURL:  getEnv("GARENA_BASE_URL", ""),
				Username: getEnv("GARENA_USERNAME", ""),
				Password: getEnv("GARENA_PASSWORD", ""),
				APIKey:   getEnv("GARENA_API_KEY", ""),
				Timeout:  getDurationEnv("GARENA_TIMEOUT", 30*time.Second),
			},
		},
		Logging: LoggingConfig{
			Level:    getEnv("LOG_LEVEL", "info"),
			Format:   getEnv("LOG_FORMAT", "json"),
			Output:   getEnv("LOG_OUTPUT", "stdout"),
			FilePath: getEnv("LOG_FILE_PATH", "logs/app.log"),
		},
		App: AppConfig{
			Name:        getEnv("APP_NAME", "Topup Application"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
			Environment: getEnv("APP_ENVIRONMENT", "development"),
			Debug:       getBoolEnv("APP_DEBUG", true),
			URL:         getEnv("APP_URL", "http://localhost:8080"),
			SecretKey:   getEnv("APP_SECRET_KEY", "your-app-secret-key"),
		},
	}

	return config, nil
}

// Utility functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// GetDSN returns the database connection string for PostgreSQL
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.Username,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
		c.Database.Timezone,
	)
}

// GetRedisAddr returns Redis connection address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

// GetServerAddr returns server address
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// IsProduction checks if app is running in production mode
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.App.Environment) == "production"
}

// IsDevelopment checks if app is running in development mode
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.App.Environment) == "development"
}