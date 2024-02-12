package webserver

import (
	"github.com/caioaraujo/go-rate-limiter/internal/infra/web/middleware"
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer() *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: ":8080",
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() {
	s.Router.Use(chimiddleware.Logger)
	for path, handler := range s.Handlers {
		s.Router.Handle(path, middleware.RateLimiter(handler))
	}
	err := http.ListenAndServe(s.WebServerPort, s.Router)
	if err != nil {
		panic(err)
	}
}
