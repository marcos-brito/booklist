package main

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/marcos-brito/booklist/internal/auth"
	"github.com/marcos-brito/booklist/internal/conn"
	"github.com/marcos-brito/booklist/internal/resolvers"
)

func setupPostgres() {
	db, err := conn.NewPostgresConnection()
	if err != nil {
		log.Fatalf("couldn't connect postgres: %s", err)
	}

	conn.InitDatabase(db)
	err = conn.Migrate(db)

	if err != nil {
		log.Fatalf("couldn't run migrations: %s", err)
	}
}

func setupRedis() {
	rdb := conn.NewRedisClient()
	err := rdb.Ping(context.Background()).Err()

	if err != nil {
		log.Fatalf("couldn't connect to redis: %s", err)
	}

	conn.InitRedis(rdb)
}

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("couldn't load .env")
	}

	ory := conn.NewOryClient()
	conn.InitOry(ory)

	setupPostgres()
	setupRedis()

	router := http.NewServeMux()
	graphql := handler.NewDefaultServer(resolvers.NewExecutableSchema(resolvers.Config{Resolvers: &resolvers.Resolver{}}))

	router.Handle("/graphql", graphql)
	router.Handle("/", playground.Handler("Booklist", "/graphql"))

	server := http.Server{
		Addr:    ":8080",
		Handler: auth.SessionMiddleware(router, conn.Ory),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("couldn't start the server: %s", err)
	}
}
