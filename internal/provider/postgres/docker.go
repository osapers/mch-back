package postgres

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/jackc/pgx/v4"

	"github.com/ory/dockertest/v3"
)

var (
	mx                  = sync.Mutex{}
	defaultDatabaseName = "mch"
)

// StartDBInDocker returns postgresql Service
func StartDBInDocker() (cleanup func() error, retURL, container string, err error) {
	var (
		repo = "bitnami/postgresql"
		tag  = "11.12.0-debian-10-r23"
		env  = []string{
			"ALLOW_EMPTY_PASSWORD=yes",
			fmt.Sprintf("POSTGRESQL_DATABASE=%s", defaultDatabaseName),
		}
		url     = "postgres://postgres@localhost:%s/%s?sslmode=disable"
		dialect = "postgres"
	)

	mx.Lock()
	defer mx.Unlock()

	pool, err := dockertest.NewPool("")
	if err != nil {
		return func() error { return nil }, "", "", fmt.Errorf("could not connect to docker: %w", err)
	}

	resource, err := pool.Run(repo, tag, env)
	if err != nil {
		return func() error { return nil }, "", "", fmt.Errorf("could not start resource: %w", err)
	}

	cleanup = func() error {
		return cleanupDockerResource(pool, resource)
	}

	url = fmt.Sprintf(url, resource.GetPort("5432/tcp"), defaultDatabaseName)

	if err := pool.Retry(func() error {
		db, err := pgx.Connect(context.Background(), url)
		if err != nil {
			return fmt.Errorf("error opening %s dev container: %w", dialect, err)
		}

		if err := db.Ping(context.Background()); err != nil {
			return err
		}
		defer func() {
			_ = db.Close(context.Background())
		}()
		return nil
	}); err != nil {
		return func() error { return nil }, "", "", fmt.Errorf("could not connect to docker: %w", err)
	}

	return cleanup, url, resource.Container.Name, nil
}

func cleanupDockerResource(pool *dockertest.Pool, resource *dockertest.Resource) error {
	var err error
	for i := 0; i < 10; i++ {
		err = pool.Purge(resource)
		if err == nil {
			return nil
		}
	}

	if err != nil && strings.Contains(err.Error(), "No such container") {
		return nil
	}

	return fmt.Errorf("failed to cleanup local container: %s", err)
}
