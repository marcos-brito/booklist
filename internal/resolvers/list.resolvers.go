package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.56

import (
	"context"
	"fmt"

	"github.com/marcos-brito/booklist/internal/auth"
	"github.com/marcos-brito/booklist/internal/models"
	"github.com/marcos-brito/booklist/internal/store"
	"gorm.io/gorm"
)

// Books is the resolver for the books field.
func (r *listResolver) Books(ctx context.Context, obj *models.List) ([]*models.Book, error) {
	books, err := store.NewListStore(store.DB).FindBooks(obj.ID)
	if err != nil {
		return nil, ErrInternal
	}

	return books, nil
}

// Owner is the resolver for the owner field.
func (r *listResolver) Owner(ctx context.Context, obj *models.List) (*models.User, error) {
	profile, err := store.NewUserStore(store.DB).FindFullProfileById(obj.ProfileID)
	if err != nil {
		return nil, ErrInternal
	}

	if profile.Settings.Private {
		return nil, nil
	}

	user := &models.User{
		UUID: profile.UUID,
	}

	return user, nil
}

// CreateList is the resolver for the createList field.
func (r *mutationResolver) CreateList(ctx context.Context, name string, description *string, publish *bool) (*models.List, error) {
	_, ident, ok := auth.GetSession(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	if publish == nil {
		publish = new(bool)
		*publish = false
	}

	list, err := store.NewListStore(store.DB).Create(name, description, *publish, ident.UUID)
	if err != nil {
		return nil, ErrInternal
	}

	return list, nil
}

// DeleteList is the resolver for the deleteList field.
func (r *mutationResolver) DeleteList(ctx context.Context, id uint) (*models.List, error) {
	_, ident, ok := auth.GetSession(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	ok, err := listIsOwned(id, ident.UUID)
	if !ok {
		return nil, err
	}

	list, err := store.NewListStore(store.DB).Delete(id)
	if err != nil {
		return nil, ErrInternal
	}

	return list, nil
}

// PublishList is the resolver for the publishList field.
func (r *mutationResolver) PublishList(ctx context.Context, id uint) (*models.List, error) {
	_, ident, ok := auth.GetSession(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	ok, err := listIsOwned(id, ident.UUID)
	if !ok {
		return nil, err
	}

	list, err := store.NewListStore(store.DB).Publish(id)
	if err != nil {
		return nil, ErrInternal
	}

	return list, nil
}

// UnpublishList is the resolver for the unpublishList field.
func (r *mutationResolver) UnpublishList(ctx context.Context, id uint) (*models.List, error) {
	_, ident, ok := auth.GetSession(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	ok, err := listIsOwned(id, ident.UUID)
	if !ok {
		return nil, err
	}

	list, err := store.NewListStore(store.DB).Unpublish(id)
	if err != nil {
		return nil, ErrInternal
	}

	return list, nil
}

// CloneList is the resolver for the cloneList field.
func (r *mutationResolver) CloneList(ctx context.Context, id uint) (*models.List, error) {
	_, ident, ok := auth.GetSession(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	listStore := store.NewListStore(store.DB)
	list, err := listStore.FindById(id)
	if err != nil {
		return nil, ErrWithOrInternal(gorm.ErrRecordNotFound, err, ErrBadId(id, "list"))
	}

	ok, _ = listIsOwned(id, ident.UUID)
	if !list.Published && !ok {
		return nil, ErrBadId(id, "list")
	}

	list, err = listStore.Clone(id, ident.UUID)
	if err != nil {
		return nil, ErrInternal
	}

	return list, nil
}

// FollowList is the resolver for the followList field.
func (r *mutationResolver) FollowList(ctx context.Context, id uint) (*models.List, error) {
	panic(fmt.Errorf("not implemented: FollowList - followList"))
}

// UnfollowList is the resolver for the unfollowList field.
func (r *mutationResolver) UnfollowList(ctx context.Context, id uint) (*models.List, error) {
	panic(fmt.Errorf("not implemented: UnfollowList - unfollowList"))
}

// AddToList is the resolver for the addToList field.
func (r *mutationResolver) AddToList(ctx context.Context, listID uint, bookID uint) (*models.List, error) {
	_, ident, ok := auth.GetSession(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	ok, err := listIsOwned(listID, ident.UUID)
	if !ok {
		return nil, err
	}

	_, err = store.NewBookStore(store.DB).FindById(bookID)
	if err != nil {
		return nil, ErrWithOrInternal(gorm.ErrRecordNotFound, err, ErrBadId(bookID, "book"))
	}

	list, err := store.NewListStore(store.DB).AddBook(listID, bookID)
	if err != nil {
		return nil, ErrInternal
	}

	return list, nil
}

// RemoveFromList is the resolver for the removeFromList field.
func (r *mutationResolver) RemoveFromList(ctx context.Context, listID uint, bookID uint) (*models.List, error) {
	_, ident, ok := auth.GetSession(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	ok, err := listIsOwned(listID, ident.UUID)
	if !ok {
		return nil, err
	}

	_, err = store.NewBookStore(store.DB).FindById(bookID)
	if err != nil {
		return nil, ErrWithOrInternal(gorm.ErrRecordNotFound, err, ErrBadId(bookID, "book"))
	}

	list, err := store.NewListStore(store.DB).RemoveBook(listID, bookID)
	if err != nil {
		return nil, ErrInternal
	}

	return list, nil
}

// List returns ListResolver implementation.
func (r *Resolver) List() ListResolver { return &listResolver{r} }

type listResolver struct{ *Resolver }