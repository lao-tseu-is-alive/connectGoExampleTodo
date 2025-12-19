package main

import (
	"log"
	"net/http"

	"connectrpc.com/connect"
	todov1 "github.com/lao-tseu-is-alive/connectGoExampleTodo/gen/todo/v1"
	"github.com/lao-tseu-is-alive/connectGoExampleTodo/gen/todo/v1/todov1connect"
	"github.com/lao-tseu-is-alive/connectGoExampleTodo/internal/todo"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	// Initialize service
	service := todo.NewTodoService()

	// ConnectRPC handler
	mux := http.NewServeMux()
	path, handler := todov1connect.NewTodoServiceHandler(service)
	mux.Handle(path, handler)

	// CORS for browser clients (e.g., gRPC-Web)
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Combined handler: Vanguard handles all traffic, proxies to Connect
	http.Handle("/", corsMiddleware.Handler(vanguardHandler))

	// Serve on HTTP/1.1 and HTTP/2 (h2c for unencrypted HTTP/2)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", h2c.NewHandler(mux, &http2.Server{})); err != nil {
		log.Fatal(err)
	}
}
