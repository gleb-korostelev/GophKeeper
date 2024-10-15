package router

import (
	"github.com/gleb-korostelev/GophKeeper/internal/middleware"
	"github.com/gleb-korostelev/GophKeeper/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func RouterInit(svc service.APIServiceI, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.GzipCompressMiddleware)
	router.Use(middleware.GzipDecompressMiddleware)
	router.Use(middleware.LoggingMiddleware(logger))
	// router.Route("/api/user", func(r chi.Router) {
	// 	r.Post("/register", svc.Register)
	// 	r.Post("/login", svc.Login)

	// 	r.Route("/", func(r chi.Router) {
	// 		r.Use(middleware.EnsureUserCookie)
	// 	})
	// })
	return router
}
