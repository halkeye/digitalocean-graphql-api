package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/digitalocean/godo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/halkeye/digitalocean-graphql-api/graph"
)

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.NewResolver()))
	// see https://gqlgen.com/reference/complexity/#custom-complexity-calculation
	h.Use(extension.FixedComplexityLimit(20)) // Dont allow complex queries

	// allow websockets
	h.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return r.Host == "example.org"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://google.com"}

	r.Use(cors.New(corsConfig))
	r.Use(DOContextToContextMiddleware())
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())
	r.Run()
}

func DOContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]string{"message": "not authorized"})
			return
		}
		client := godo.NewFromToken(bearerToken)
		ctx := context.WithValue(c.Request.Context(), graph.DOContextKey, client)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	token := strings.Split(header, " ")
	if len(token) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return token[1], nil
}
