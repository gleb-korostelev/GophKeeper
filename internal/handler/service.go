package handler

import (
	"net/http"
)

type API interface {
	// Healthcheck
	Healthcheck(rw http.ResponseWriter, r *http.Request)
}

type ProfileSvc interface {
}

type Implementation struct {
	ProfileSvc ProfileSvc
}

func NewImplementation(profileSvc ProfileSvc) API {
	return &Implementation{
		ProfileSvc: profileSvc,
	}
}
