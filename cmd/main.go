package cmd

import (
	"log"
	"net/http"
	"os"

	accounthttp "github.com/FrostBitzX/smart-task-ai/cmd/account/http"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/config"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/database"

	"github.com/getkin/kin-openapi/openapi3"
	chi "github.com/go-chi/chi/v5"
	chi_middleware "github.com/oapi-codegen/nethttp-middleware"
)

func main() {
	cfg := config.NewConfig()
	_ = database.NewDB(cfg)

	// Load OpenAPI specification
	loader := &openapi3.Loader{Context: nil}
	spec, err := loader.LoadFromFile("openapi/openapi.yml")
	if err != nil {
		log.Fatalf("failed to load OpenAPI spec: %v", err)
	}

	if err := spec.Validate(loader.Context); err != nil {
		log.Fatalf("invalid OpenAPI spec: %v", err)
	}

	// Create chi router and attach validator middleware
	r := chi.NewRouter()
	r.Use(chi_middleware.OapiRequestValidator(spec))

	// Initialize domain services and HTTP handlers
	// TODO: Wire db into account service when repository layer is implemented
	accountHandler := accounthttp.NewAccountHandler(nil)

	strictHandler := accounthttp.NewStrictHandler(accountHandler, nil)

	// Mount generated routes under base router
	r.Mount("/", accounthttp.HandlerFromMux(strictHandler, r))

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	log.Printf("starting HTTP server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
