package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	appmw "github.com/vikramsinghpanwar/ludo-backend/internal/app/http/middleware"
	//since chi and app both have middle wares we use aliases here
	"github.com/vikramsinghpanwar/ludo-backend/internal/auth"
	"github.com/vikramsinghpanwar/ludo-backend/internal/player"
)

type Dependencies struct {
	AuthHandler   *auth.Handler
	PlayerHandler *player.Handler
}

func NewRouter(deps *Dependencies) http.Handler {
	r := chi.NewRouter()

	// global middlewares
	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// health check
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	// auth routes
	if deps.AuthHandler != nil {
		r.Route("/auth", func(r chi.Router) {
			auth.RegisterRoutes(r, deps.AuthHandler)
		})
	}

	// player routes
	if deps.PlayerHandler != nil {
		r.Route("/player", func(r chi.Router) {
			r.Use(appmw.Auth) // 🔐 protected
			player.RegisterRoutes(r, deps.PlayerHandler)
		})
	}

	return r
}
