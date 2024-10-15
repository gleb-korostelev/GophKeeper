package router

import (
	"net/http"

	"github.com/gleb-korostelev/GophKeeper/internal/handler"
	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	"github.com/gleb-korostelev/GophKeeper/tools/swagger"
	"github.com/gorilla/mux"
)

const AppName = "gophkeeper"

// CreateRouter
func CreateRouter(impl handler.API, appPort int, isSwaggerCreated bool) *mux.Router {
	// mw := middleware.NewCoreMW()
	var handlers = []swagger.Handler{
		{
			HandlerFunc:      impl.Healthcheck,
			Path:             "/healthcheck",
			Method:           http.MethodGet,
			Description:      "Healthcheck",
			ResponseBody:     response.HealthcheckResp{},
			ResponseMimeType: swagger.MimeJson,
			Opts:             []swagger.Option{},
			Tag:              Hc,
		},
	}

	return NewAPI(AppName, appPort, isSwaggerCreated, handlers)
}
