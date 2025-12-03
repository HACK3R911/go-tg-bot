package handler

import (
	"github.com/HACK3R911/go-tg-bot/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}
