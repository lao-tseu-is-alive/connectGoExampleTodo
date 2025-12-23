package main

import (
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"connectrpc.com/vanguard"
	"github.com/lao-tseu-is-alive/connectGoExampleTodo/gen/todo/v1/todov1connect"
	"github.com/lao-tseu-is-alive/connectGoExampleTodo/pkg/todo"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
)

const (
	APP     = "todoServer"
	VERSION = "0.0.1"
	//AppVersion      = APP + " v" + VERSION
	defaultIp      = "127.0.0.1"
	defaultPort    = 8080
	defaultLogName = "stderr"
	myHeaderKey    = "Acme-Tenant-Id"
	//serverHeaderKey = "App-Version"
)

func main() {

	logWriter, err := config.GetLogWriter(defaultLogName)
	if err != nil {
		log.Fatalf("💥💥 error getting log writer: %v'\n", err)
	}
	logLevel, err := config.GetLogLevel(golog.InfoLevel)
	if err != nil {
		log.Fatalf("💥💥 error getting log level: %v'\n", err)
	}
	l := golog.NewLogger("simple", logWriter, logLevel, APP)
	l.Info("🚀 Starting", "app", APP, "version", VERSION)
	// Initialize service
	todoService := todo.NewTodoService(myHeaderKey, l)

	// Create the Connect handler with validation interceptor
	path, handler := todov1connect.NewTodoServiceHandler(
		todoService,
		// Validation via Protovalidate is almost always recommended
		connect.WithInterceptors(validate.NewInterceptor()),
	)

	// Wrap with Vanguard for REST transcoding support
	// This enables REST-style URLs like GET /v1/greet/{name} and POST /v1/greet
	service := vanguard.NewService(path, handler)
	transcoder, err := vanguard.NewTranscoder([]*vanguard.Service{service})
	if err != nil {
		log.Fatalf("💥💥 error creating vanguard transcoder: %v'\n", err)
	}

	p := new(http.Protocols)
	p.SetHTTP1(true)
	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)
	serverAddress := fmt.Sprintf("%s:%d", defaultIp, defaultPort)
	l.Info("Will start server ...", "listenAddress", serverAddress)
	l.Info("REST endpoints available:", "GET", "/v1/todos", "POST", "/v1/greet")
	s := http.Server{
		Addr:      serverAddress,
		Handler:   transcoder, // Use transcoder instead of mux
		Protocols: p,
	}
	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("💥💥 error doing ListenAndServe: %v'\n", err)
	}
	/*
		// CORS for browser clients (e.g., gRPC-Web)
		corsMiddleware := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"*"},
		})


		// Vanguard transcoder (for REST)
		vanguardHandler, err := vanguard.NewTranscoder([]*vanguard.Service{
			vanguard.NewService(todov1connect.TodoServiceName, connect.NewClient[todov1.TodoServiceClient](
				http.DefaultClient,
				"http://localhost:8080", // Point to internal Connect server
				connect.WithGRPC(),
			)),
		})
		if err != nil {
			log.Fatal(err)
		}
		// Combined handler: Vanguard handles all traffic, proxies to Connect
		http.Handle("/", corsMiddleware.Handler(vanguardHandler))

		// Serve on HTTP/1.1 and HTTP/2 (h2c for unencrypted HTTP/2)
		log.Println("Server starting on :8080")
		if err := http.ListenAndServe(":8080", h2c.NewHandler(mux, &http2.Server{})); err != nil {
			log.Fatal(err)
		}

	*/
}
