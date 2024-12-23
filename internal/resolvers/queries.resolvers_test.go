package resolvers_test

import (
	"context"
	"testing"

	"github.com/marcos-brito/booklist/internal/auth"
	"github.com/marcos-brito/booklist/internal/models"
	"github.com/marcos-brito/booklist/internal/resolvers"
	ory "github.com/ory/client-go"
)

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
