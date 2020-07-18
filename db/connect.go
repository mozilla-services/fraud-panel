package db // import "go.mozilla.org/fraud-panel/db"

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"time"

	// lib/pq is the postgres driver
	_ "github.com/lib/pq"

	"go.mozilla.org/fraud-panel/mozlog"
)

// Handler handles a database connection
type Handler struct {
	*sql.DB
}

// Transaction owns a sql transaction
type Transaction struct {
	*sql.Tx
	ID uint64
}

// Config holds the parameters to connect to a database
type Config struct {
	Name                string
	User                string
	Password            string
	Host                string
	SSLMode             string
	SSLRootCert         string
	MaxOpenConns        int
	MaxIdleConns        int
	MonitorPollInterval time.Duration
}

var h *sql.DB

// Connect creates a database connection and returns a handler
func Connect(config Config) error {
	var (
		dsn string
		err error
	)
	if os.Getenv("DB_DSN") != "" {
		dsn = os.Getenv("DB_DSN")
	} else {
		userPass := url.UserPassword(config.User, config.Password)
		if config.SSLMode == "" {
			config.SSLMode = "disable"
		}
		dsn = fmt.Sprintf("postgres://%s@%s/%s?sslmode=%s&sslrootcert=%s",
			userPass.String(), config.Host, config.Name, config.SSLMode, config.SSLRootCert)
	}
	h, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	if config.MaxOpenConns > 0 {
		h.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		h.SetMaxIdleConns(config.MaxIdleConns)
	}
	dbCheckCtx, dbCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dbCancel()
	err = CheckConnectionContext(dbCheckCtx)
	if err != nil {
		return err
	}
	mozlog.Info("database connection established", nil)
	return nil
}

// CheckConnectionContext runs a test query against the database and
// returns an error if it fails
func CheckConnectionContext(ctx context.Context) error {
	var one uint
	err := h.QueryRowContext(ctx, "SELECT 1").Scan(&one)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	if one != 1 {
		return fmt.Errorf("failed connection check: `select 1` returned %d", one)
	}
	return nil
}

// Monitor queries the database every pollInterval and kills the program
// if it becomes unavailable
func Monitor(pollInterval time.Duration) {
	mozlog.Info("starting DB monitor polling every "+pollInterval.String(), nil)
	for {
		err := CheckConnectionContext(context.Background())
		if err != nil {
			mozlog.Info(err.Error(), nil)
			os.Exit(10)
		}
		time.Sleep(pollInterval)
	}
}
