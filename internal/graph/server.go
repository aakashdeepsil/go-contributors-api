package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aakashdeepsil/go-contributors-api/internal/graph/generated"
	"github.com/aakashdeepsil/go-contributors-api/internal/graph/resolvers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(resolver *resolvers.Resolver) *chi.Mux {
	router := chi.NewRouter()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	return router
}
