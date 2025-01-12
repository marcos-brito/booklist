package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.56

import (
	"context"

	"github.com/marcos-brito/booklist/internal/auth"
	"github.com/marcos-brito/booklist/internal/conn"
	"github.com/marcos-brito/booklist/internal/models"
	"github.com/marcos-brito/booklist/internal/store"
)

// Settings is the resolver for the settings field.
func (r *currentUserResolver) Settings(ctx context.Context, obj *models.CurrentUser) (*models.Settings, error) {
	settings, err := store.NewUserStore(conn.DB).FindSettingsByUserUuid(obj.UUID)
	if err != nil {
		return nil, ErrInternal
	}

	return settings, nil
}

// Lists is the resolver for the lists field.
func (r *currentUserResolver) Lists(ctx context.Context, obj *models.CurrentUser) ([]*models.List, error) {
	lists, err := store.NewUserStore(conn.DB).FindLists(obj.UUID)
	if err != nil {
		return nil, ErrInternal
	}

	return lists, nil
}

// Collection is the resolver for the collection field.
func (r *currentUserResolver) Collection(ctx context.Context, obj *models.CurrentUser) ([]*models.CollectionItem, error) {
	items, err := store.NewUserStore(conn.DB).FindItems(obj.UUID)
	if err != nil {
		return nil, ErrInternal
	}

	return items, nil
}

// UpdateSettings is the resolver for the updateSettings field.
func (r *mutationResolver) UpdateSettings(ctx context.Context, changes models.UpdateSettings) (*models.Settings, error) {
	_, ident, ok := auth.GetSession(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	settings, err := store.NewUserStore(conn.DB).UpdateSettings(ident.UUID, changes)
	if err != nil {
		return nil, ErrInternal
	}

	return settings, nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*models.CurrentUser, error) {
	_, ident, ok := auth.GetSession(ctx)
	if !ok {
		return nil, nil
	}

	_, err := store.NewUserStore(conn.DB).FindProfileByUserUuid(ident.UUID)
	if err != nil {
		return nil, ErrInternal
	}

	user := &models.CurrentUser{
		UUID:  ident.UUID,
		Name:  ident.Traits.Name,
		Email: ident.Traits.Email,
	}

	return user, nil
}

// CurrentUser returns CurrentUserResolver implementation.
func (r *Resolver) CurrentUser() CurrentUserResolver { return &currentUserResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type currentUserResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
