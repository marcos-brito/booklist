package resolvers_test

import (
	"context"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/marcos-brito/booklist/internal/auth"
	"github.com/marcos-brito/booklist/internal/models"
	"github.com/marcos-brito/booklist/internal/resolvers"
	"github.com/marcos-brito/booklist/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
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
	err = store.Migrate(db)
	if err != nil {
		panic(err)
	}

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

func TestSettings(t *testing.T) {
	resolver := &resolvers.Resolver{}
	ctx, profile := NewUser(t)
	_, err := resolver.CurrentUser().Settings(ctx, profile)

	assert.Nil(t, err)
}

func TestLists(t *testing.T) {
	resolver := &resolvers.Resolver{}

	t.Run("should return user lists", func(t *testing.T) {
		ctx, user := NewUser(t)
		created := []*models.List{
			CreateList(t, ctx, false),
			CreateList(t, ctx, true),
		}

		lists, err := resolver.CurrentUser().Lists(ctx, user)
		assert.Nil(t, err)
		for _, list := range created {
			assert.Contains(t, lists, list)
		}
	})
}

func TestCollection(t *testing.T) {
	resolver := &resolvers.Resolver{}
	ctx, _ := NewUser(t)
	inputs := []models.CreateBook{
		{
			Title: "O homem de giz",
			Isbn:  "9788551002933",
		},
		{
			Title: "Arquitetura limpa",
			Isbn:  "9788550804606",
		},
	}

	books := []*models.Book{}
	for _, input := range inputs {
		book, err := resolver.Mutation().CreateBook(ctx, input)
		assert.Nil(t, err)
		books = append(books, book)
	}

	t.Run("should return the user's collection", func(t *testing.T) {
		ctx, user := NewUser(t)

		for _, book := range books {
			AddItemToUserCollection(t, ctx, book.ID)
		}

		got, err := resolver.CurrentUser().Collection(ctx, user)
		assert.Nil(t, err)
		for _, item := range got {
			assert.True(t, slices.ContainsFunc(books, func(book *models.Book) bool {
				return book.ID == item.BookID
			}))
		}
	})
}

func TestUpdateSettings(t *testing.T) {
	resolver := resolvers.Resolver{}

	t.Run("should changes user settings", func(t *testing.T) {
		ctx, user := NewUser(t)
		changes := models.UpdateSettings{
			ShowName:       true,
			ShowCollection: true,
		}

		_, err := resolver.Mutation().UpdateSettings(ctx, changes)
		assert.Nil(t, err)

		got, err := resolver.CurrentUser().Settings(ctx, user)
		assert.Nil(t, err)
		assert.True(t, got.ShowName)
		assert.True(t, got.ShowCollection)
	})

	t.Run("should fail if there is no session", func(t *testing.T) {
		ctx, _ := NewUser(t)
		ctx = auth.AddSessionToContext(ctx, nil)
		settings, err := resolver.Mutation().UpdateSettings(ctx, models.UpdateSettings{})

		assert.Nil(t, settings)
		assert.ErrorIs(t, err, resolvers.ErrUnauthorized)
	})
}

func TestMe(t *testing.T) {
	resolver := resolvers.Resolver{}

	t.Run("should return current user", func(t *testing.T) {
		session := NewRandomSession()
		ctx := auth.AddSessionToContext(context.Background(), session)

		_, ident, ok := auth.GetSession(ctx)
		assert.True(t, ok)

		got, err := resolver.Query().Me(ctx)
		assert.Nil(t, err)
		assert.Equal(t, got.UUID, ident.UUID)
		assert.Equal(t, got.Name, ident.Traits.Name)
		assert.Equal(t, got.Email, ident.Traits.Email)
	})

	t.Run("should return nil if there is no session", func(t *testing.T) {
		ctx := auth.AddSessionToContext(context.Background(), nil)
		got, err := resolver.Query().Me(ctx)

		assert.Nil(t, err)
		assert.Nil(t, got)
	})
}
