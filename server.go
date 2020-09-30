package main

import (
	"QuicPos/graph"
	"QuicPos/graph/generated"
	"QuicPos/internal/mongodb"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

//server
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mongodb.InitDB()
	defer mongodb.DisconnectDB()
	mongodb.List()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
