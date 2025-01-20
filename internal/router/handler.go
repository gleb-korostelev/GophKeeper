// Package router provides utilities for creating and managing HTTP routers,
// including support for middleware and Swagger documentation generation.
package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gleb-korostelev/GophKeeper/middleware"
	"github.com/gleb-korostelev/GophKeeper/tools/logger"
	"github.com/gleb-korostelev/GophKeeper/tools/swagger"
	"github.com/gorilla/mux"
	v3 "github.com/swaggest/swgui/v3"
)

// Swagger represents a configuration structure for generating Swagger documentation.
//
// Fields:
// - nameAPI: The name of the API for Swagger documentation.
// - handlers: A slice of `swagger.Handler` objects representing API endpoints.
// - JSON: The Swagger JSON documentation content.
type Swagger struct {
	nameAPI  string
	handlers []swagger.Handler
	JSON     string
}

// NewAPI initializes a new HTTP router and optionally generates Swagger documentation.
//
// Parameters:
// - nameAPI: The name of the API, used in Swagger documentation.
// - appPort: The main application port. Swagger documentation will be served on `appPort+1`.
// - isSwaggerCreated: A boolean indicating whether Swagger documentation has already been created.
// - handlers: A slice of `swagger.Handler` objects defining API endpoints.
//
// Returns:
// - *mux.Router: A configured Gorilla Mux router with registered handlers and middleware.
//
// Workflow:
// 1. Sets up middleware using `middleware.PanicMid`.
// 2. Registers the provided handlers to the router.
// 3. Generates Swagger documentation if not already created.
// 4. Serves the Swagger JSON documentation and UI on `appPort+1`.
//
// Example usage:
//
//	handlers := []swagger.Handler{
//	    {Path: "/api/v1/example", Method: "GET", HandlerFunc: exampleHandler},
//	}
//	router := NewAPI("ExampleAPI", 8080, false, handlers)
//
// Swagger Documentation:
// - The Swagger JSON file is generated and saved to `./internal/router/json/swagger.json`.
// - The Swagger UI is accessible at `/docs`.
//
// Middleware:
// - Adds `middleware.PanicMid` to handle panics and unexpected errors gracefully.
func NewAPI(nameAPI string, appPort int, isSwaggerCreated bool, handlers []swagger.Handler) *mux.Router {
	// Return nil if no handlers are provided.
	if len(handlers) == 0 {
		return nil
	}

	// Initialize a new Gorilla Mux router.
	r := mux.NewRouter()

	// Apply panic handling middleware.
	r.Use(middleware.PanicMid)

	// Register handlers to the router.
	for _, h := range handlers {
		r.HandleFunc(h.Path, h.HandlerFunc).Methods(h.Method)
	}

	// Generate Swagger documentation if not already created.
	if !isSwaggerCreated {
		swaggerDoc, err := swagger.GenerateDoc(nameAPI, handlers)
		if err != nil {
			logger.Fatal(err)
		}

		// Save the generated Swagger documentation to a file.
		err = os.WriteFile("./internal/router/json/swagger.json", []byte(swaggerDoc), os.ModePerm)
		if err != nil {
			logger.Fatal(err)
		}
	}

	// Serve Swagger JSON and UI.
	fs := http.FileServer(http.Dir("./internal/router/json"))
	http.Handle("/swagger.json", fs)
	http.Handle("/", v3.NewHandler(nameAPI, "/swagger.json", "/docs"))

	// Start a separate HTTP server for Swagger documentation on `appPort+1`.
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", appPort+1), nil)
	}()

	return r
}
