package resolvers_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/marcos-brito/booklist/internal/store"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {
    teardown := setup()
	code := m.Run()

    teardown()
	os.Exit(code)
}

func setup() func() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	container := startPostgres()
	db, err := store.NewConnection()
	if err != nil {
		panic(err)
	}

	store.With(db)
	store.Migrate(db)

	return func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			panic(fmt.Sprintf("failed to terminate container: %s", err))
		}
	}
}

func startPostgres() testcontainers.Container {
    ctx := context.Background()
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(os.Getenv("POSTGRES_DB")),
		postgres.WithUsername(os.Getenv("POSTGRES_USER")),
		postgres.WithPassword(os.Getenv("POSTGRES_PASSWORD")),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		panic(fmt.Sprintf("failed to start container: %s", err))
	}

    port,err  := container.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}

    os.Setenv("POSTGRES_PORT", port.Port())

	return container
}
