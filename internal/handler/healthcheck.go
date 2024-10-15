package handler

import (
	"net/http"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
)

func (i *Implementation) Healthcheck(rw http.ResponseWriter, r *http.Request) {
	response.Healthcheck(rw)
}
