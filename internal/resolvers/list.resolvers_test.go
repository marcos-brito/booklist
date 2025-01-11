package resolvers_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/marcos-brito/booklist/internal/auth"
	"github.com/marcos-brito/booklist/internal/models"
	"github.com/marcos-brito/booklist/internal/resolvers"
	"github.com/stretchr/testify/assert"
)

func CreateList(t *testing.T, ctx context.Context, publish bool) *models.List {
	resolver := resolvers.Resolver{}
	list, err := resolver.Mutation().CreateList(ctx, fmt.Sprintf("List %d", rand.Int()), nil, &publish)

	assert.Nil(t, err)
	return list
}

func UpdateSettings(t *testing.T, ctx context.Context, settings *models.UpdateSettings) {
	resolver := resolvers.Resolver{}
	_, err := resolver.Mutation().UpdateSettings(ctx, *settings)

	assert.Nil(t, err)
}

func TestBooks(t *testing.T) {
}

func TestOwner(t *testing.T) {
	resolver := resolvers.Resolver{}
	ctx, owner := NewUser(t)
	list := CreateList(t, ctx, true)

	t.Run("should return list owner", func(t *testing.T) {
		UpdateSettings(t, ctx, &models.UpdateSettings{Private: false})
		ctx, _ := NewUser(t)
		got, err := resolver.List().Owner(ctx, list)

		assert.Nil(t, err)
		assert.Equal(t, owner.UUID, got.UUID)
	})

	t.Run("should return nil if user profile is private", func(t *testing.T) {
		UpdateSettings(t, ctx, &models.UpdateSettings{Private: true})
		ctx, _ := NewUser(t)
		got, err := resolver.List().Owner(ctx, list)
		assert.Nil(t, err)
		assert.Nil(t, got)
	})
}

func TestCreateList(t *testing.T) {
	resolver := resolvers.Resolver{}

	t.Run("should create private list by default", func(t *testing.T) {
		ctx, user := NewUser(t)

		list, err := resolver.Mutation().CreateList(ctx, "list", nil, nil)
		assert.Nil(t, err)
		assert.False(t, list.Published)

		lists, err := resolver.CurrentUser().Lists(ctx, user)
		assert.Nil(t, err)
		assert.Contains(t, lists, list)
	})

	t.Run("should create public list if specified", func(t *testing.T) {
		ctx, user := NewUser(t)
		publish := true

		list, err := resolver.Mutation().CreateList(ctx, "list", nil, &publish)
		assert.Nil(t, err)
		assert.True(t, list.Published)

		lists, err := resolver.CurrentUser().Lists(ctx, user)
		assert.Nil(t, err)
		assert.Contains(t, lists, list)
	})

	t.Run("should fail it there is no session", func(t *testing.T) {
		ctx, _ := NewUser(t)
		ctx = auth.AddSessionToContext(ctx, nil)
		got, err := resolver.Mutation().CreateList(ctx, "list", nil, nil)

		assert.Nil(t, got)
		assert.ErrorIs(t, err, resolvers.ErrUnauthorized)
	})
}

func TestDeleteList(t *testing.T) {
	resolver := resolvers.Resolver{}

	t.Run("should delete a list", func(t *testing.T) {
		ctx, user := NewUser(t)
		list := CreateList(t, ctx, false)

		deleted, err := resolver.Mutation().DeleteList(ctx, list.ID)
		assert.Nil(t, err)

		lists, err := resolver.CurrentUser().Lists(ctx, user)
		assert.Nil(t, err)
		assert.NotContains(t, lists, deleted)
	})

	t.Run("should fail it there is no session", func(t *testing.T) {
		ctx, _ := NewUser(t)
		list := CreateList(t, ctx, false)
		ctx = auth.AddSessionToContext(ctx, nil)
		got, err := resolver.Mutation().DeleteList(ctx, list.ID)

		assert.Nil(t, got)
		assert.ErrorIs(t, err, resolvers.ErrUnauthorized)
	})

	t.Run("should fail if user does not own list", func(t *testing.T) {
		ctx1, _ := NewUser(t)
		ctx2, _ := NewUser(t)
		list := CreateList(t, ctx1, false)
		got, err := resolver.Mutation().DeleteList(ctx2, list.ID)

		assert.Nil(t, got)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadId(list.ID, "list")))
	})
}

func TestPublishList(t *testing.T) {
	resolver := resolvers.Resolver{}

	t.Run("should publish a list", func(t *testing.T) {
		ctx, _ := NewUser(t)
		list := CreateList(t, ctx, false)
		list, err := resolver.Mutation().PublishList(ctx, list.ID)

		assert.Nil(t, err)
		assert.True(t, list.Published)
	})

	t.Run("shoud fail if there is no session", func(t *testing.T) {
		ctx, _ := NewUser(t)
		list := CreateList(t, ctx, false)
		ctx = auth.AddSessionToContext(ctx, nil)
		list, err := resolver.Mutation().PublishList(ctx, list.ID)

		assert.Nil(t, list)
		assert.ErrorIs(t, err, resolvers.ErrUnauthorized)
	})

	t.Run("should fail if user does not own the list", func(t *testing.T) {
		ctx1, _ := NewUser(t)
		ctx2, _ := NewUser(t)
		list := CreateList(t, ctx1, false)
		got, err := resolver.Mutation().PublishList(ctx2, list.ID)

		assert.Nil(t, got)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadId(list.ID, "list")))
	})
}

func TestUnpublishList(t *testing.T) {
	resolver := resolvers.Resolver{}

	t.Run("should unpublish a list", func(t *testing.T) {
		ctx, _ := NewUser(t)
		list := CreateList(t, ctx, true)
		list, err := resolver.Mutation().UnpublishList(ctx, list.ID)

		assert.Nil(t, err)
		assert.False(t, list.Published)
	})

	t.Run("shoud fail if there is no session", func(t *testing.T) {
		ctx, _ := NewUser(t)
		list := CreateList(t, ctx, true)
		ctx = auth.AddSessionToContext(ctx, nil)
		list, err := resolver.Mutation().UnpublishList(ctx, list.ID)

		assert.Nil(t, list)
		assert.ErrorIs(t, err, resolvers.ErrUnauthorized)
	})

	t.Run("should fail if user does not own the list", func(t *testing.T) {
		ctx1, _ := NewUser(t)
		ctx2, _ := NewUser(t)
		list := CreateList(t, ctx1, true)
		got, err := resolver.Mutation().UnpublishList(ctx2, list.ID)

		assert.Nil(t, got)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadId(list.ID, "list")))
	})
}

func TestCloneList(t *testing.T) {
	resolver := resolvers.Resolver{}
	ctx, user := NewUser(t)
	published := CreateList(t, ctx, true)
	private := CreateList(t, ctx, false)

	t.Run("should clone published list", func(t *testing.T) {
		ctx, user := NewUser(t)

		got, err := resolver.Mutation().CloneList(ctx, published.ID)
		assert.Nil(t, err)

		lists, err := resolver.CurrentUser().Lists(ctx, user)
		assert.Nil(t, err)
		assert.Equal(t, lists[0].ID, got.ID)
	})

	t.Run("should clone private list when it's owned", func(t *testing.T) {
		got, err := resolver.Mutation().CloneList(ctx, private.ID)
		assert.Nil(t, err)

		lists, err := resolver.CurrentUser().Lists(ctx, user)
		assert.Nil(t, err)
		assert.Equal(t, lists[2].ID, got.ID)
	})

	t.Run("should fail if list is not published", func(t *testing.T) {
		ctx, _ := NewUser(t)
		list, err := resolver.Mutation().CloneList(ctx, private.ID)

		assert.Nil(t, list)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadId(private.ID, "list")))
	})

	t.Run("should fail if there is no session", func(t *testing.T) {
		ctx, _ := NewUser(t)
		ctx = auth.AddSessionToContext(ctx, nil)
		list, err := resolver.Mutation().CloneList(ctx, published.ID)

		assert.Nil(t, list)
		assert.ErrorIs(t, err, resolvers.ErrUnauthorized)
	})
}

func TestFollowList(t *testing.T) {}

func TestUnfollowList(t *testing.T) {}

func TestAddToList(t *testing.T) {
	resolver := resolvers.Resolver{}
	ctx, _ := NewUser(t)
	book := CreateBook(t, ctx)

	t.Run("should add a book to a list", func(t *testing.T) {
		ctx, _ := NewUser(t)
		list := CreateList(t, ctx, false)

		list, err := resolver.Mutation().AddToList(ctx, list.ID, book.ID)
		assert.Nil(t, err)

		books, err := resolver.List().Books(ctx, list)
		assert.Nil(t, err)
		assert.Equal(t, books[0].ID, book.ID)
	})

	t.Run("should fail if user does not own the list", func(t *testing.T) {
		ctx1, _ := NewUser(t)
		ctx2, _ := NewUser(t)
		list := CreateList(t, ctx1, false)
		got, err := resolver.Mutation().AddToList(ctx2, list.ID, book.ID)

		assert.Nil(t, got)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadId(list.ID, "list")))
	})

	t.Run("should fail if book id does not exist", func(t *testing.T) {
		ctx, _ := NewUser(t)
		list := CreateList(t, ctx, false)
		id := rand.Uint()
		list, err := resolver.Mutation().AddToList(ctx, list.ID, id)

		assert.Nil(t, list)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadId(id, "book")))
	})

	t.Run("should fail if there is no session", func(t *testing.T) {
		ctx, _ := NewUser(t)
		list := CreateList(t, ctx, false)
		ctx = auth.AddSessionToContext(ctx, nil)
		list, err := resolver.Mutation().AddToList(ctx, list.ID, book.ID)

		assert.Nil(t, list)
		assert.ErrorIs(t, err, resolvers.ErrUnauthorized)
	})
}

func TestDeleteFromList(t *testing.T) {
	resolver := resolvers.Resolver{}
	ctx, _ := NewUser(t)
	book := CreateBook(t, ctx)

	t.Run("should delete a book from a list", func(t *testing.T) {
		ctx, _ := NewUser(t)
		list := CreateList(t, ctx, false)

		list, err := resolver.Mutation().AddToList(ctx, list.ID, book.ID)
		assert.Nil(t, err)

		list, err = resolver.Mutation().RemoveFromList(ctx, list.ID, book.ID)
		assert.Nil(t, err)

		books, err := resolver.List().Books(ctx, list)
		assert.Nil(t, err)
		assert.Empty(t, books)
	})

	t.Run("should fail if user does not own the list", func(t *testing.T) {
		ctx1, _ := NewUser(t)
		ctx2, _ := NewUser(t)
		list := CreateList(t, ctx1, false)
		got, err := resolver.Mutation().RemoveFromList(ctx2, list.ID, book.ID)

		assert.Nil(t, got)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadId(list.ID, "list")))
	})

	t.Run("should fail if book id does not exist", func(t *testing.T) {
		ctx, _ := NewUser(t)
		list := CreateList(t, ctx, false)
		id := rand.Uint()
		list, err := resolver.Mutation().RemoveFromList(ctx, list.ID, id)

		assert.Nil(t, list)
		assert.ErrorAs(t, err, ErrAsPointer(resolvers.ErrBadId(id, "book")))
	})

	t.Run("should fail if there is no session", func(t *testing.T) {
		ctx, _ := NewUser(t)
		list := CreateList(t, ctx, false)
		ctx = auth.AddSessionToContext(ctx, nil)
		list, err := resolver.Mutation().RemoveFromList(ctx, list.ID, book.ID)

		assert.Nil(t, list)
		assert.ErrorIs(t, err, resolvers.ErrUnauthorized)
	})
}
