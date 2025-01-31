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
type Swagger struct {
	nameAPI  string
	handlers []swagger.Handler
	JSON     string
}

// NewAPI initializes a new HTTP router and optionally generates Swagger documentation.
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
