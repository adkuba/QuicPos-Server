package main

import (
	"QuicPos/graph"
	"QuicPos/graph/generated"
	"QuicPos/internal/ip"
	"QuicPos/internal/mongodb"
	"QuicPos/internal/storage"
	"QuicPos/internal/tensorflow"
	"QuicPos/internal/user"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

const defaultPort = "8080"

//server
func main() {
	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		//Debug:            true,
		AllowedHeaders: []string{"Access-Control-Allow-Origin", "Content-Type"},
	}).Handler)

	router.Use(ip.Middleware())

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	storage.InitStorage()
	tensorflow.InitModels()

	mongodb.InitDB()
	defer mongodb.DisconnectDB()

	user.CheckCounter()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
