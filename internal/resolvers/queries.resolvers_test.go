package resolvers_test

import (
	"context"
	"testing"
	"os"
	"time"
    "fmt"

	"github.com/marcos-brito/booklist/internal/auth"
	"github.com/marcos-brito/booklist/internal/models"
	"github.com/marcos-brito/booklist/internal/resolvers"
	"github.com/joho/godotenv"
	"github.com/marcos-brito/booklist/internal/store"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	ory "github.com/ory/client-go"
)

func TestMain(m *testing.M) {
	teardown := Setup()
    defer teardown()
	code := m.Run()
	os.Exit(code)
}

func Setup() func() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	container := StartPostgres()
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

func StartPostgres() testcontainers.Container {
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

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}

	os.Setenv("POSTGRES_PORT", port.Port())

	return container
}

func TestMeResolver(t *testing.T) {
	tests := []struct {
		desc     string
		session  *ory.Session
		expected *models.Profile
	}{
		{
			"creates profile with default settings",
			&ory.Session{
                Identity: &ory.Identity{Id: "123" ,Traits: &auth.Identity{Name: "User", Email: "user@email.com"}}},
			&models.Profile{
				UserUUID: "123",
				Name:     "User",
				Email:    "user@email.com",
				Settings: models.Settings{
					Private: true,
				},
			},
		},
		{
			"returns nil if there is no session",
			nil,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			ctx := auth.AddSessionToContext(context.Background(), test.session)
			resolver := resolvers.Resolver{}
			got, err := resolver.Query().Me(ctx)

			if err != nil {
				t.Fatalf("unexpected error %s", err)
			}

			if !profilePartialEquals(got, test.expected) {
				t.Fatalf("got mismatching profiles:\n got\n %+v\n expected\n %+v", got, test.expected)
			}
		})
	}
}

func profilePartialEquals(l, r *models.Profile) bool {
    if (l == nil && r == nil) {
        return true
    }

    return l.UUID() == r.UUID() && l.Email == r.Email && l.Name == r.Name && l.Settings.Private == r.Settings.Private
}
