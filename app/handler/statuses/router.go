package statuses

import (
	"net/http"

	"yatter-backend-go/app/app"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi"
)

// Implementation of handler
type handler struct {
	app *app.App
}

// Create Handler for `/v1/statuses/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()
	h := &handler{app: app}

	// DOC: https://go-chi.io/#/pages/routing?id=routing-groups
	r.Route("/", func(r chi.Router) {
		// 以下の処理は認証を必要とする
		r.Use(auth.Middleware(app))
		r.Post("/", h.Create)
	})

	r.Get("/{id}", h.Get)

	return r
}
