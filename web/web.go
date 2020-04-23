package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/polisgo2020/search-senyast4745/index"
	"github.com/polisgo2020/search-senyast4745/util"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/polisgo2020/search-senyast4745/config"
	"github.com/rs/zerolog/log"
)

type App struct {
	Mux          *chi.Mux
	ind          Indexed
	netInterface string
}

type FileResponse struct {
	Filename string
	Count    int
	Spacing  int
}

type Indexed interface {
	GetIndex(str ...string) (*index.Index, error)
}

func NewApp(c *config.Config, i Indexed) (*App, error) {
	r := chi.NewMux()

	log.Debug().Msg("add custom log and header middleware")

	r.Use(logMiddleware)
	r.Use(headerMiddleware)

	d, err := time.ParseDuration(c.TimeOut)
	if err != nil {
		log.Warn().Str("timeout", c.TimeOut).Msg("can not parse timeout")
		d = 10 * time.Millisecond
	}

	log.Debug().Dur("timeout", d).Msg("server timeout")

	r.Use(middleware.Timeout(d))

	corsFilter := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	log.Debug().Interface("cors policy", corsFilter).Msg("cors policy created")

	r.Use(corsFilter.Handler)

	log.Debug().RawJSON("endpoint", []byte("{\"method\" : \"POST\", \"pattern\" : \"\\\"")).Msg("register controller")

	app := &App{Mux: r, netInterface: c.Listen, ind: i}

	r.Post("/", app.searchHandler)
	return app, nil
}

func (a *App) searchHandler(w http.ResponseWriter, req *http.Request) {
	searchWords := req.FormValue("search")
	log.Info().Str("search phrase", searchWords).Msg("start search")
	var inputWords []string
	for _, word := range strings.Split(searchWords, " ") {
		util.CleanUserInput(word, func(input string) {
			inputWords = append(inputWords, input)
		})
	}
	log.Debug().Msgf("clean input: %+v", inputWords)
	if len(inputWords) == 0 {
		log.Err(nil).Str("input", searchWords).Msg("Incorrect search words")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var resp []FileResponse
	ind, err := a.ind.GetIndex(inputWords...)
	if err != nil {
		log.Err(err).Msg("error while getting index")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for k, v := range ind.Search(inputWords) {
		resp = append(resp, FileResponse{
			Filename: k,
			Count:    v.Path,
			Spacing:  v.Weight,
		})
	}
	log.Info().Interface("result", resp).Msgf("search finished")
	log.Debug().Msg("start marshalling and writing data to response")
	rawData, err := json.Marshal(resp)
	if err != nil {
		log.Err(err).Interface("json data", resp).Msg("error while marshalling data to json")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if _, err = fmt.Fprint(w, string(rawData)); err != nil {
		log.Printf("error %s while writing data %s do json\n", err, string(rawData))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	log.Debug().Interface("headers", w.Header())
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		log.Debug().
			Str("method", r.Method).
			Str("remote", r.RemoteAddr).
			Str("path", r.URL.Path).
			Int("duration", int(time.Since(start))).
			Msgf("Called url %s", r.URL.Path)
	})

}

func (a *App) Run() {
	log.Info().Str("network interface", a.netInterface).Msg("server start")
	if err := http.ListenAndServe(a.netInterface, a.Mux); err != nil {
		log.Err(err).Str("network interface", a.netInterface).Msg("can't start server")
	}
	log.Info().Str("network interface", a.netInterface).Msg("server shutdown")
}
