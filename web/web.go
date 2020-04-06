package web

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/polisgo2020/search-senyast4745/config"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

type App struct {
	Mux          *chi.Mux
	netInterface string
}

func NewApp(c *config.Config, middlewares ...func(handler http.Handler) http.Handler) (*App, error) {
	r := chi.NewMux()
	for i := range middlewares {
		r.Use(middlewares[i])
	}
	d, err := time.ParseDuration(c.TimeOut)
	if err != nil {
		d = 10
	}

	r.Use(middleware.Timeout(d * time.Millisecond))
	filesDir := http.Dir("static")
	corsFilter := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(corsFilter.Handler)
	fileServer(r, "/", filesDir)
	if c.Listen == "" {
		c.Listen = ":8888"
	}
	return &App{Mux: r, netInterface: c.Listen}, nil
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func (a *App) Run() {
	log.Debug().Msg("server shutdown")
	if err := http.ListenAndServe(a.netInterface, a.Mux); err != nil {
		log.Err(err).Str("network interface", a.netInterface).Msg("can't start server")
	}
}
