package resolvers_test

import (
	"math/rand/v2"
	"testing"

	"github.com/marcos-brito/booklist/internal/auth"
	"github.com/marcos-brito/booklist/internal/models"
	"github.com/marcos-brito/booklist/internal/resolvers"
	"github.com/stretchr/testify/assert"
)

func TestCreateBook(t *testing.T) {
	resolver := resolvers.Resolver{}

	t.Run("should allow create with missing fields", func(t *testing.T) {
		ctx, _ := NewUser(t)
		input := models.CreateBook{
			Title: "A maldição da casa das flores",
			Isbn:  "9786555664973",
		}

		got, err := resolver.Mutation().CreateBook(ctx, input)
		assert.Nil(t, err)

        // addedBy, err := resolver.Book().AddedBy(ctx, got)
		assert.Nil(t, err)
		assert.True(t, got.NeedsApproval)
		assert.Equal(t, got.ISBN, input.Isbn)
		assert.Equal(t, got.Title, input.Title)
		// assert.Equal(t, addedBy.UUID, user.UUID)
	})

	t.Run("should fail if there is no session", func(t *testing.T) {
		ctx, _ := NewUser(t)
		input := models.CreateBook{
			Title: "A maldição da casa das flores",
			Isbn:  "9786555664973",
		}

		ctx = auth.AddSessionToContext(ctx, nil)
		got, err := resolver.Mutation().CreateBook(ctx, input)
		assert.ErrorIs(t, err, resolvers.ErrUnauthorized)
		assert.Nil(t, got)
	})

	t.Run("should fail if publisher id does not exist", func(t *testing.T) {
		ctx, _ := NewUser(t)
		publisherId := uint(rand.Uint32())
		input := models.CreateBook{
			Title:     "A maldição da casa das flores",
			Isbn:      "9786555664973",
			Publisher: &publisherId,
		}

		got, err := resolver.Mutation().CreateBook(ctx, input)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadId(publisherId, "publisher")))
		assert.Nil(t, got)
	})

	t.Run("should fail if any author id does not exist", func(t *testing.T) {
		ctx, _ := NewUser(t)
		input := models.CreateBook{
			Title:   "A maldição da casa das flores",
			Isbn:    "9786555664973",
			Authors: []uint{uint(rand.Uint32()), uint(rand.Uint32()), uint(rand.Uint32())},
		}

		got, err := resolver.Mutation().CreateBook(ctx, input)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadId(input.Authors[0], "author")))
		assert.Nil(t, got)
	})
}
