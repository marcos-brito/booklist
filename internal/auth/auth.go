package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	ory "github.com/ory/client-go"
)

type SessionContextKey string

type Identity struct {
	UUID   uuid.UUID
	Traits *Traits
	*ory.Identity
}

type Traits struct {
	Name  string
	Email string
}

const session_context_key SessionContextKey = "req.session"

func newIdentityFromOry(oryIdent *ory.Identity) (*Identity, error) {
	traits, err := ParseTraits(oryIdent.Traits)
	if err != nil {
		return nil, err
	}

	uuid, err := uuid.Parse(oryIdent.Id)
	if err != nil {
		return nil, err
	}

	ident := &Identity{
		UUID:     uuid,
		Traits:   traits,
		Identity: oryIdent,
	}

	return ident, nil
}

func ParseTraits(obj interface{}) (*Traits, error) {
	id, ok := obj.(*Traits)

	if !ok {
		return nil, fmt.Errorf("can't parse %v into identity", obj)
	}

	return id, nil
}

func SessionMiddleware(next http.Handler, oryClient *ory.APIClient) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		cookies := request.Header.Get("Cookie")
		session, _, _ := oryClient.FrontendAPI.ToSession(request.Context()).Cookie(cookies).Execute()
		ctx := AddSessionToContext(request.Context(), session)

		next.ServeHTTP(writer, request.WithContext(ctx))
	}
}

func AddSessionToContext(ctx context.Context, session *ory.Session) context.Context {
	return context.WithValue(ctx, session_context_key, session)
}

func GetSession(ctx context.Context) (*ory.Session, *Identity, bool) {
	session, ok := ctx.Value(session_context_key).(*ory.Session)
	if !ok || session == nil {
		return nil, nil, false
	}

	ident, err := newIdentityFromOry(session.Identity)
	if err != nil {
		return nil, nil, false
	}

	return session, ident, true
}

func FindIdentity(uuid uuid.UUID, client *ory.APIClient) (*Identity, bool) {
	resp, _, err := client.IdentityAPI.GetIdentity(context.Background(), uuid.String()).Execute()
	if err != nil {
		return nil, false
	}

	ident, err := newIdentityFromOry(resp)
	if err != nil {
		return nil, false
	}

	return ident, true
}
