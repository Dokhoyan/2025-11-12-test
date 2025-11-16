package handler

import (
	"github.com/Dokhoyan/2025-11-12-test/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}
