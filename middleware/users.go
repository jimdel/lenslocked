package middleware

import (
	"net/http"

	"github.com/jimdel/lenslocked/context"
	"github.com/jimdel/lenslocked/controllers"
	"github.com/jimdel/lenslocked/models"
)

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (um *UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := controllers.ReadCookie(r, controllers.CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := um.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithUser(r.Context(), user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
