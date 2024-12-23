package main

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/joho/godotenv/autoload"
	"github.com/marcos-brito/booklist/internal/auth"
	"github.com/marcos-brito/booklist/internal/resolvers"
	"github.com/marcos-brito/booklist/internal/store"
	ory "github.com/ory/client-go"
)

var oryClient *ory.APIClient

func init() {
	config := ory.NewConfiguration()
	oryClient = ory.NewAPIClient(config)
}

func main() {
	db, err := store.NewConnection()
	if err != nil {
		panic(fmt.Errorf("couldn't start the server: %s", err))
	}

	store.With(db)
	store.Migrate(db)

	router := http.NewServeMux()
	graphql := handler.NewDefaultServer(resolvers.NewExecutableSchema(resolvers.Config{Resolvers: &resolvers.Resolver{}}))

	router.Handle("/graphql", graphql)
	router.Handle("/", playground.Handler("Booklist", "/graphql"))

	server := http.Server{
		Addr:    ":8080",
		Handler: auth.SessionMiddleware(router, oryClient),
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Errorf("couldn't start the server: %s", err))
	}
}
