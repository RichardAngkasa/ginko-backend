package handlers

import (
	"net/http"
)

func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
			MaxAge:   -1,
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("logged out"))
	}
}
