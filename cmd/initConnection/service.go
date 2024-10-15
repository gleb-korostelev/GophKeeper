package initConnection

import (
	"context"
	"net/http"

	"github.com/gleb-korostelev/GophKeeper/config"
	"github.com/gleb-korostelev/GophKeeper/internal/handler"
	"github.com/gleb-korostelev/GophKeeper/internal/router"
	"github.com/gleb-korostelev/GophKeeper/service/profile"
	"github.com/gleb-korostelev/GophKeeper/tools/db"
	"github.com/rs/cors"
)

func InitImpl(
	ctx context.Context,
	adapter db.IAdapter,
	port int,
) http.Handler {

	isSwaggerCreated := config.GetConfigBool(config.IsSwaggerCreated)

	profileSvc := initServices(adapter)

	api := handler.NewImplementation(profileSvc)
	r := router.CreateRouter(api, port, isSwaggerCreated)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodDelete,
			http.MethodPost,
			http.MethodPut,
			http.MethodOptions,
		},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Content-Length", "Content-Type", "Access-Control-Allow-Origin"},
	})

	return c.Handler(r)
}

// initServices func. Here you have to implement business logic or repository services
func initServices(db db.IAdapter) (
	profileSvc handler.ProfileSvc,
) {
	profileSvc = profile.NewService(db)

	return
}
