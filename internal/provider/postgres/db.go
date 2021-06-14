// Package postgres provides DB functionality
package postgres

//go:generate go-bindata -pkg migrations -ignore bindata -prefix ./migrations/ -o ./migrations/bindata.go ./migrations

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/osapers/mch-back/internal/provider/postgres/migrations"
	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // nolint:golint
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

// Conn contains postgres required settings and connection
type Conn struct {
	DB     *pg.DB
	Secret []byte
	logger *zap.Logger
}

// New returns Conn instance
func New() (*Conn, error) {
	cfg, err := newConfig()
	if err != nil {
		return nil, err
	}

	connOptions, err := pg.ParseURL(cfg.connectionURL)
	if err != nil {
		return nil, fmt.Errorf("cannot parse connectionURL - %s", err)
	}

	if err = applyMigrations(cfg.connectionURL); err != nil {
		return nil, err
	}

	l := newLogger()

	connOptions.OnConnect = func(_ context.Context, _ *pg.Conn) error {
		l.Info("successfully connected to PostgreSQL")
		return nil
	}

	db, err := open(connOptions)
	if err != nil {
		return nil, err
	}

	db.AddQueryHook(debugHook{
		logger:  l,
		verbose: cfg.logVerbose,
	})

	conn := &Conn{
		DB:     db,
		Secret: []byte(cfg.secret),
		logger: l,
	}

	return conn, nil
}

// Destroy releases the connection and sync db logger
func (c *Conn) Destroy() error {
	_ = c.logger.Sync()

	if err := c.DB.Close(); err != nil {
		return fmt.Errorf("cannot close db connection properly - %s", err)
	}

	return nil
}

// open a database connection which is long-lived.
// You need to call Destroy() on the returned pgx.Conn
func open(connOpts *pg.Options) (*pg.DB, error) {
	if connOpts == nil {
		return nil, errors.New("invalid connection db options")
	}

	db := pg.Connect(connOpts)

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("cannot connect to postgres - %s", err)
	}

	return db, nil
}

// Migrate a database schema
func applyMigrations(connectionURL string) error {
	if connectionURL == "" {
		return errors.New("db connectionURL is unset")
	}

	sourceInstance, err := bindata.WithInstance(bindata.Resource(migrations.AssetNames(), migrations.Asset))
	if err != nil {
		return fmt.Errorf("unable to create source of migrations: %w", err)
	}

	// run migrations
	m, err := migrate.NewWithSourceInstance("go-bindata", sourceInstance, connectionURL)
	if err != nil {
		return fmt.Errorf("unable to create migrations: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("unable to run migrations: %w", err)
	}

	return nil
}
