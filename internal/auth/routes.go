package auth

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Post("/signup", h.Signup)
	r.Post("/otp", h.OTP)
	r.Post("/verify-otp", h.VerifyOTP)
	r.Post("/login", h.Login)
	r.Post("/logout", h.Logout)
	r.Post("/refresh", h.Refresh)
}
