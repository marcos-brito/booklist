package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	ory "github.com/ory/client-go"
)

type Identity struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ParseIdentity(obj interface{}) (*Identity, error) {
	id, ok := obj.(*Identity)

	if !ok {
		return nil, fmt.Errorf("can't parse %v into identity", obj)
	}

	return id, nil
}

const session_context_key string = "req.session"

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

func GetSession(ctx context.Context) (*ory.Session, error) {
	session, ok := ctx.Value(session_context_key).(*ory.Session)

	if !ok || session == nil {
		return nil, errors.New("session not found in context")
	}

	return session, nil
}
