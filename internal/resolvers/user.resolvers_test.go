package resolvers_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/marcos-brito/booklist/internal/resolvers"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	resolver := resolvers.Resolver{}
	_, user := NewUser(t)

	t.Run("shold return the user", func(t *testing.T) {
		got, err := resolver.Query().User(context.Background(), user.UUID)
		assert.Nil(t, got)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadUuid(user.UUID, "user")))
	})

	t.Run("should fail if profile is private", func(t *testing.T) {
		got, err := resolver.Query().User(context.Background(), user.UUID)
		assert.Nil(t, got)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadUuid(user.UUID, "user")))
	})

	t.Run("should fail if uuid does not exist", func(t *testing.T) {
		uuid, err := uuid.NewUUID()
		assert.Nil(t, err)

		got, err := resolver.Query().User(context.Background(), uuid)
		assert.Nil(t, got)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadUuid(uuid, "user")))
	})
}
