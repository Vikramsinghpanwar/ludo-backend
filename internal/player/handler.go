package player

import (
	"fmt"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {

	fmt.Println("GetProfile called")
}
