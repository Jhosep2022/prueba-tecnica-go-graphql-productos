package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/Jhosep2022/prueba-tecnica-go-graphql-productos/graph/generated"
	graphqldelivery "github.com/Jhosep2022/prueba-tecnica-go-graphql-productos/internal/delivery/graphql"
	"github.com/Jhosep2022/prueba-tecnica-go-graphql-productos/internal/repository/memory"
	"github.com/Jhosep2022/prueba-tecnica-go-graphql-productos/internal/usecase"
)

const defaultPort = "8080"

func main() {
	productRepository := memory.NewProductRepository()
	productUseCase := usecase.NewProductUseCase(productRepository)
	resolver := graphqldelivery.NewResolver(productUseCase)

	graphQLServer := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{Resolvers: resolver},
		),
	)
	graphQLServer.SetErrorPresenter(graphqldelivery.ErrorPresenter)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("Product GraphQL Playground", "/query"))
	mux.Handle("/query", graphQLServer)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("GraphQL playground available at http://localhost:%s/", port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server stopped unexpectedly: %v", err)
	}
}
