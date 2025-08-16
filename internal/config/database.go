// internal/config/database.go
package config

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// DatabaseConnection holds database connections
type DatabaseConnection struct {
	SqlDB  *sql.DB
	Config *DatabaseConfig
}

// NewDatabaseConnection creates a new database connection
func NewDatabaseConnection(config *DatabaseConfig) (*DatabaseConnection, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.DBName,
		config.SSLMode,
		config.Timezone,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DatabaseConnection{
		SqlDB:  db,
		Config: config,
	}, nil
}

// Close closes the database connection
func (dc *DatabaseConnection) Close() error {
	if dc.SqlDB != nil {
		return dc.SqlDB.Close()
	}
	return nil
}

// GetDB returns the sql.DB instance
func (dc *DatabaseConnection) GetDB() *sql.DB {
	return dc.SqlDB
}

// Health checks database health
func (dc *DatabaseConnection) Health(ctx context.Context) error {
	return dc.SqlDB.PingContext(ctx)
}

// GetConnectionStats returns database connection statistics
func (dc *DatabaseConnection) GetConnectionStats() sql.DBStats {
	return dc.SqlDB.Stats()
}

// BeginTransaction starts a new database transaction
func (dc *DatabaseConnection) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return dc.SqlDB.BeginTx(ctx, nil)
}

// BeginTransactionWithOptions starts a transaction with specific isolation level
func (dc *DatabaseConnection) BeginTransactionWithOptions(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return dc.SqlDB.BeginTx(ctx, opts)
}

// WithTransaction executes a function within a database transaction
func (dc *DatabaseConnection) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := dc.SqlDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithTransactionRetry executes a function within a transaction with retry logic
func (dc *DatabaseConnection) WithTransactionRetry(ctx context.Context, maxRetries int, fn func(tx *sql.Tx) error) error {
	var lastErr error
	
	for i := 0; i < maxRetries; i++ {
		err := dc.WithTransaction(ctx, fn)
		if err == nil {
			return nil
		}
		
		lastErr = err
		
		// Check if error is retryable (serialization failure, deadlock, etc.)
		if !isRetryableError(err) {
			return err
		}
		
		// Wait before retrying with exponential backoff
		waitTime := time.Duration(i+1) * 100 * time.Millisecond
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
		}
	}
	
	return fmt.Errorf("transaction failed after %d retries: %w", maxRetries, lastErr)
}

// isRetryableError checks if database error is retryable
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := strings.ToLower(err.Error())
	retryableErrors := []string{
		"serialization failure",
		"deadlock detected",
		"connection reset",
		"connection refused",
		"timeout",
	}
	
	for _, retryable := range retryableErrors {
		if strings.Contains(errStr, retryable) {
			return true
		}
	}
	
	return false
}

// PrepareStatement prepares a SQL statement for repeated use
func (dc *DatabaseConnection) PrepareStatement(ctx context.Context, query string) (*sql.Stmt, error) {
	stmt, err := dc.SqlDB.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	return stmt, nil
}

// ExecContext executes a query without returning any rows
func (dc *DatabaseConnection) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return dc.SqlDB.ExecContext(ctx, query, args...)
}

// QueryContext executes a query that returns rows
func (dc *DatabaseConnection) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return dc.SqlDB.QueryContext(ctx, query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row
func (dc *DatabaseConnection) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return dc.SqlDB.QueryRowContext(ctx, query, args...)
}

// Migration functions
func (dc *DatabaseConnection) CreateDatabase(ctx context.Context, dbName string) error {
	// Create a connection without specifying database name
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s sslmode=%s",
		dc.Config.Host,
		dc.Config.Port,
		dc.Config.Username,
		dc.Config.Password,
		dc.Config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect for database creation: %w", err)
	}
	defer db.Close()

	// Check if database exists
	var exists bool
	query := "SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)"
	err = db.QueryRowContext(ctx, query, dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	// Create database if it doesn't exist
	if !exists {
		// Note: Database names cannot be parameterized, so we need to validate the name
		if !isValidDatabaseName(dbName) {
			return fmt.Errorf("invalid database name: %s", dbName)
		}
		
		createQuery := fmt.Sprintf("CREATE DATABASE %s", dbName)
		_, err = db.ExecContext(ctx, createQuery)
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}

	return nil
}

func (dc *DatabaseConnection) DropDatabase(ctx context.Context, dbName string) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s sslmode=%s",
		dc.Config.Host,
		dc.Config.Port,
		dc.Config.Username,
		dc.Config.Password,
		dc.Config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect for database deletion: %w", err)
	}
	defer db.Close()

	if !isValidDatabaseName(dbName) {
		return fmt.Errorf("invalid database name: %s", dbName)
	}

	dropQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
	_, err = db.ExecContext(ctx, dropQuery)
	if err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	return nil
}

func isValidDatabaseName(name string) bool {
	if len(name) == 0 || len(name) > 63 {
		return false
	}
	

	for i, r := range name {
		if i == 0 {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_') {
				return false
			}
		} else {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || 
				 (r >= '0' && r <= '9') || r == '_' || r == '$') {
				return false
			}
		}
	}
	
	return true
}

func (dc *DatabaseConnection) RunMigration(ctx context.Context, migrationSQL string) error {
	return dc.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, migrationSQL)
		return err
	})
}

func (dc *DatabaseConnection) RunMigrationFile(ctx context.Context, filePath string) error {
	
	return fmt.Errorf("not implemented - read file %s and call RunMigration", filePath)
}

// DatabaseHealthCheck represents database health status
type DatabaseHealthCheck struct {
	Status      string        `json:"status"`
	Message     string        `json:"message"`
	Latency     time.Duration `json:"latency"`
	Connections sql.DBStats   `json:"connections"`
}

// HealthCheck performs a comprehensive database health check
func (dc *DatabaseConnection) HealthCheck(ctx context.Context) *DatabaseHealthCheck {
	start := time.Now()
	
	healthCheck := &DatabaseHealthCheck{
		Connections: dc.GetConnectionStats(),
	}

	// Test database connectivity
	if err := dc.Health(ctx); err != nil {
		healthCheck.Status = "unhealthy"
		healthCheck.Message = fmt.Sprintf("Database connection failed: %v", err)
		healthCheck.Latency = time.Since(start)
		return healthCheck
	}

	healthCheck.Status = "healthy"
	healthCheck.Message = "Database is accessible"
	healthCheck.Latency = time.Since(start)

	return healthCheck
}

// Bulk operations for better performance
type BulkInsertResult struct {
	RowsAffected int64
	LastInsertID int64
	Error        error
}

// BulkInsert performs bulk insert operation
func (dc *DatabaseConnection) BulkInsert(ctx context.Context, tableName string, columns []string, values [][]interface{}) *BulkInsertResult {
	if len(values) == 0 {
		return &BulkInsertResult{Error: fmt.Errorf("no values to insert")}
	}

	// Build bulk insert query
	placeholders := make([]string, len(values))
	args := make([]interface{}, 0, len(values)*len(columns))
	
	for i, row := range values {
		if len(row) != len(columns) {
			return &BulkInsertResult{Error: fmt.Errorf("row %d has %d values, expected %d", i, len(row), len(columns))}
		}
		
		rowPlaceholders := make([]string, len(columns))
		for j, val := range row {
			argIndex := i*len(columns) + j + 1
			rowPlaceholders[j] = fmt.Sprintf("$%d", argIndex)
			args = append(args, val)
		}
		placeholders[i] = fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ","))
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		tableName,
		strings.Join(columns, ","),
		strings.Join(placeholders, ","),
	)

	result, err := dc.SqlDB.ExecContext(ctx, query, args...)
	if err != nil {
		return &BulkInsertResult{Error: err}
	}

	rowsAffected, _ := result.RowsAffected()
	lastInsertID, _ := result.LastInsertId()

	return &BulkInsertResult{
		RowsAffected: rowsAffected,
		LastInsertID: lastInsertID,
	}
}

// GetTableColumns returns column information for a table
func (dc *DatabaseConnection) GetTableColumns(ctx context.Context, tableName string) ([]string, error) {
	query := `
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = $1 
		ORDER BY ordinal_position
	`
	
	rows, err := dc.SqlDB.QueryContext(ctx, query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns: %w", err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}
		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return columns, nil
}

func (dc *DatabaseConnection) TableExists(ctx context.Context, tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = $1
		)
	`
	
	var exists bool
	err := dc.SqlDB.QueryRowContext(ctx, query, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check table existence: %w", err)
	}

	return exists, nil
}