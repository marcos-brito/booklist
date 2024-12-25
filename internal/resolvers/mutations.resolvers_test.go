package resolvers_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/marcos-brito/booklist/internal/auth"
	"github.com/marcos-brito/booklist/internal/models"
	"github.com/marcos-brito/booklist/internal/resolvers"
	ory "github.com/ory/client-go"
	"github.com/stretchr/testify/assert"
)

func NewUser(t *testing.T) (context.Context, *models.Profile) {
	resolver := resolvers.Resolver{}
	ctx := auth.AddSessionToContext(context.Background(), NewRandomSession())

	profile, err := resolver.Query().Me(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	return ctx, profile
}

func NewRandomSession() *ory.Session {
	return &ory.Session{Identity: &ory.Identity{
		Id: uuid.NewString(),
		Traits: &auth.Identity{
			Name:  "user",
			Email: fmt.Sprintf("user%d@email.com", rand.Int()),
		}}}
}

func AddItemToUserCollection(t *testing.T, ctx context.Context, bookId uint) *models.CollectionItem {
	resolver := resolvers.Resolver{}
	status := models.StatusReading
	item, err := resolver.Mutation().AddToCollection(ctx, bookId, &status)

	assert.Nil(t, err)
	return item
}

// Take a error, usually a pointer, and returns
// a pointer to it. With this is possible
// to assert inline.
func ErrAsPointer(err error) *error {
	return &err
}

func TestCreateBook(t *testing.T) {
	resolver := resolvers.Resolver{}

	t.Run("should allow create with missing fields", func(t *testing.T) {
		ctx, profile := NewUser(t)
		input := models.CreateBook{
			Title: "A maldição da casa das flores",
			Isbn:  "9786555664973",
		}

		got, err := resolver.Mutation().CreateBook(ctx, input)
		assert.Nil(t, err)
		assert.True(t, got.NeedsApproval)
		assert.Equal(t, got.ISBN, input.Isbn)
		assert.Equal(t, got.Title, input.Title)
		assert.Equal(t, got.ProfileID, profile.ID)
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

func TestAddToCollection(t *testing.T) {
	resolver := resolvers.Resolver{}
	ctx, _ := NewUser(t)
	book, err := resolver.Mutation().CreateBook(ctx, models.CreateBook{
		Title: "Não conta a ninguém",
		Isbn:  "9788580419689",
	})

	assert.Nil(t, err)

	t.Run("should add item to collection", func(t *testing.T) {
		ctx, profile := NewUser(t)
		status := models.StatusDropped

		created, err := resolver.Mutation().AddToCollection(ctx, book.ID, &status)
		assert.Nil(t, err)

		collection, err := resolver.Profile().Collection(ctx, profile)
		assert.Nil(t, err)
		assert.True(t, slices.ContainsFunc(collection, func(item *models.CollectionItem) bool {
			return item.ID == created.ID
		}))
	})

	t.Run("should fail if book id does not exist", func(t *testing.T) {
		ctx, _ := NewUser(t)
		id := uint(rand.Uint32())
		status := models.StatusDropped
		got, err := resolver.Mutation().AddToCollection(ctx, id, &status)

		assert.Nil(t, got)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadId(id, "book")))
	})

	t.Run("should fail if there is no session", func(t *testing.T) {
		ctx, _ := NewUser(t)
		id := uint(rand.Uint32())
		status := models.StatusDropped
		ctx = auth.AddSessionToContext(ctx, nil)
		got, err := resolver.Mutation().AddToCollection(ctx, id, &status)

		assert.Nil(t, got)
		assert.ErrorIs(t, err, resolvers.ErrUnauthorized)
	})

}

func TestDeleteFromCollection(t *testing.T) {
	resolver := resolvers.Resolver{}
	ctx, _ := NewUser(t)
	book, err := resolver.Mutation().CreateBook(ctx, models.CreateBook{
		Title: "Seis anos depois",
		Isbn:  "9786555650310",
	})

	assert.Nil(t, err)

	t.Run("should delete item from collection", func(t *testing.T) {
		ctx, profile := NewUser(t)
		item := AddItemToUserCollection(t, ctx, book.ID)

		deleted, err := resolver.Mutation().DeleteFromCollection(ctx, item.ID)
		assert.Nil(t, err)

		collection, err := resolver.Profile().Collection(ctx, profile)
		assert.Nil(t, err)
		assert.False(t, slices.ContainsFunc(collection, func(item *models.CollectionItem) bool {
			return item.ID == deleted.ID
		}))
	})

	t.Run("should fail if item id does not exist", func(t *testing.T) {
		ctx, _ := NewUser(t)
		id := uint(rand.Uint32())
		got, err := resolver.Mutation().DeleteFromCollection(ctx, id)

		assert.Nil(t, got)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadId(id, "collectionItem")))
	})

	t.Run("should fail if there is no session", func(t *testing.T) {
		ctx, _ := NewUser(t)
		id := uint(rand.Uint32())
		ctx = auth.AddSessionToContext(ctx, nil)
		got, err := resolver.Mutation().DeleteFromCollection(ctx, id)

		assert.Nil(t, got)
		assert.ErrorIs(t, err, resolvers.ErrUnauthorized)
	})

}

func TestChangeItemStatus(t *testing.T) {
	resolver := resolvers.Resolver{}
	ctx, _ := NewUser(t)
	book, err := resolver.Mutation().CreateBook(ctx, models.CreateBook{
		Title: "Contra todas as probabilidades do amor",
		Isbn:  "9788595810105",
	})

	assert.Nil(t, err)

	t.Run("should change item status", func(t *testing.T) {
		ctx, profile := NewUser(t)
		item := AddItemToUserCollection(t, ctx, book.ID)

		changed, err := resolver.Mutation().ChangeItemStatus(ctx, item.ID, models.StatusDropped)
		assert.Nil(t, err)
		assert.Equal(t, changed.Status, models.StatusDropped)

		collection, err := resolver.Profile().Collection(ctx, profile)
		assert.Nil(t, err)
		assert.True(t, slices.ContainsFunc(collection, func(item *models.CollectionItem) bool {
			return item.ID == changed.ID && changed.Status == models.StatusDropped
		}))
	})

	t.Run("should fail it user does not own item id", func(t *testing.T) {
		ctx1, _ := NewUser(t)
		ctx2, _ := NewUser(t)
		item := AddItemToUserCollection(t, ctx1, book.ID)
		got, err := resolver.Mutation().ChangeItemStatus(ctx2, item.ID, models.StatusOnHold)

		assert.Nil(t, got)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadId(item.ID, "collectionItem")))
	})

	t.Run("should fail if item id does not exist", func(t *testing.T) {
		ctx, _ := NewUser(t)
		id := uint(rand.Uint32())
		got, err := resolver.Mutation().ChangeItemStatus(ctx, id, models.StatusOnHold)

		assert.Nil(t, got)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadId(id, "collectionItem")))
	})

	t.Run("should fail there is no session", func(t *testing.T) {
		ctx, _ := NewUser(t)
		item := AddItemToUserCollection(t, ctx, book.ID)
		ctx = auth.AddSessionToContext(ctx, nil)
		got, err := resolver.Mutation().ChangeItemStatus(ctx, item.ID, models.StatusOnHold)

		assert.Nil(t, got)
		assert.ErrorIs(t, err, resolvers.ErrUnauthorized)
	})
}
