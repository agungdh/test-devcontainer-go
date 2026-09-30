package docs

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	apdocs "todoapp/docs"
)

const scalarPage = `<!doctype html>
<html>
<head>
  <title>Todo API Docs</title>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
</head>
<body>
  <script id="api-reference" data-url="/docs/openapi.yaml"></script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(scalarPage))
}

func (h *Handler) Spec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	_, _ = w.Write(apdocs.Spec)
}

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/docs", func(r chi.Router) {
		r.Get("/", h.Index)
		r.Get("/openapi.yaml", h.Spec)
	})
}
