package core

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/migrate"
)

const (
	dbHostEnv     string = "DB_HOST"
	dbPortEnv     string = "DB_PORT"
	dbUserEnv     string = "DB_USER"
	dbPasswordEnv string = "DB_PASSWORD"
	dbNameEnv     string = "DB_NAME"
	dbSslModeEnv  string = "DB_SSLMODE"

	defaultSslMode string = "require"
)

type DbConfig struct {
	host     string
	port     string
	user     string
	password string
	name     string
	sslMode  string
}

func newDbConfig() (*DbConfig, error) {
	host := os.Getenv(dbHostEnv)
	if strings.TrimSpace(host) == "" {
		return nil, fmt.Errorf("missing environment variable: %s", dbHostEnv)
	}

	port := os.Getenv(dbPortEnv)
	if strings.TrimSpace(port) == "" {
		return nil, fmt.Errorf("Missing environment variable: %s", dbPortEnv)
	}

	user := os.Getenv(dbUserEnv)
	if strings.TrimSpace(user) == "" {
		return nil, fmt.Errorf("Missing environment variable: %s", dbUserEnv)
	}

	password := os.Getenv(dbPasswordEnv)
	if strings.TrimSpace(password) == "" {
		return nil, fmt.Errorf("Missing environment variable: %s", dbPasswordEnv)
	}

	name := os.Getenv(dbNameEnv)
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("Missing environment variable: %s", dbNameEnv)
	}

	sslMode := os.Getenv(dbSslModeEnv)
	if strings.TrimSpace(sslMode) == "" {
		log.Printf("Missing environment variable for %s, defaulting to %s", dbSslModeEnv, defaultSslMode)
		sslMode = defaultSslMode
	}

	return &DbConfig{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		name:     name,
		sslMode:  sslMode,
	}, nil
}

func (config *DbConfig) ConnectionString() string {
	// TODO: determine whether sslmode needs to be enabled for prod or not
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.user, config.password, config.host, config.port, config.name, config.sslMode)
}

func runMigrations(ctx context.Context, db *bun.DB, migrationFiles fs.FS) error {
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(migrationFiles); err != nil {
		return fmt.Errorf("Failed to discover migrations: %w", err)
	}
	migrator := migrate.NewMigrator(db, migrations)
	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("Failed to init migrator: %w", err)
	}
	group, err := migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("Failed to run migrations: %w", err)
	}

	if group.IsZero() {
		log.Println("No new migrations to run")
	} else {
		log.Println("Migrated to", group)
	}

	return nil
}

// SetupDb initializes and configures a database connection using environment variables.
// It creates a connection pool, runs migrations, and returns a *bun.DB instance.
//
// Parameters:
//   - ctx: A context.Context for managing the database operations.
//   - migrationFiles: An fs.FS containing the migration files to be applied.
//
// Returns:
//   - *bun.DB: A pointer to the configured bun.DB instance.
//   - error: An error if any step in the setup process fails.
//
// The function performs the following steps:
// 1. Creates a database configuration from environment variables.
// 2. Parses the configuration to create a connection pool.
// 3. Opens a SQL database from the pool.
// 4. Initializes a bun.DB instance with the SQL database.
// 5. Pings the database to ensure connectivity.
// 6. Runs database migrations using the provided migration files.
//
// If any step fails, it returns an error with a descriptive message.
func SetupDb(ctx context.Context, migrationFiles fs.FS) (*bun.DB, error) {
	dbConfig, err := newDbConfig()
	if err != nil {
		return nil, fmt.Errorf("Could not create db config: %w", err)
	}
	pgxConfig, err := pgxpool.ParseConfig(dbConfig.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("Could not create pgx pool config from db config: %w", err)
	}
	pgxPool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not create pgx pool: %w", err)
	}
	sqlDb := stdlib.OpenDBFromPool(pgxPool)
	db := bun.NewDB(sqlDb, pgdialect.New())

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("Unable to ping db: %w", err)
	}

	if err := runMigrations(ctx, db, migrationFiles); err != nil {
		return nil, err
	}

	return db, nil
}
