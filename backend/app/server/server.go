package server

import (
 "fmt"
 //"log"
 "time"
 "net/http"
 "strings"

 log "github.com/go-pkgz/lgr"
 um "github.com/go-pkgz/rest"
 "github.com/pkg/errors"

 "github.com/didip/tollbooth/v6"
 "github.com/didip/tollbooth_chi"
 "github.com/go-chi/chi/v5"
 "github.com/go-chi/chi/v5/middleware"

 secret "manager-secrets/backend/app/store"
)

type Server struct {
    DataStore      secret.Store
	PinSize        int
	MaxPinAttempts int
	WebRoot        string
	Version        string
}


func (s Server) Run() error {
	log.Printf("[INFO] activate rest server")
	if err := http.ListenAndServe(":8080", s.routes()); err != http.ErrServerClosed {
		//return errors.Wrap(err, "server failed")
		return errors.Wrap(err, "server failed")
	}

	return nil
}

func (s Server) routes() chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.RequestID, middleware.RealIP, um.Recoverer(log.Default()))
	router.Use(middleware.Throttle(1000), middleware.Timeout(60*time.Second))
	router.Use(um.AppInfo("secrets", "jtrw", s.Version), um.Ping, um.SizeLimit(64*1024))
	router.Use(tollbooth_chi.LimitHandler(tollbooth.NewLimiter(10, nil)))

    router.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
          //fmt.Fprintf(w, s.DataStore.StorePath)
        //  fmt.Fprintf(w, "\n")
          fmt.Fprintf(w, s.DataStore.Get("test/secret", "one"))
           //fmt.Fprintf(w, "Secret:". s.DataStore.Get("test/secret", "one"))
    })

	router.Route("/api/v1", func(r chi.Router) {
	    r.Get("/*", s.getValuesByKey)
		//r.Use(Logger(log.Default()))
		//r.Post("/message", s.saveMessageCtrl)
		//r.Get("/message/{key}/{pin}", s.getMessageCtrl)
		//r.Get("/params", s.getParamsCtrl)
	})
//
// 	router.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
// 		render.PlainText(w, r, "User-agent: *\nDisallow: /api/\nDisallow: /show/\n")
// 	})
//
// 	s.fileServer(router, "/", http.Dir(s.WebRoot))
	return router
}

func (s Server) getValuesByKey(w http.ResponseWriter, r *http.Request) {
    key := chi.URLParam(r, "*")

    chunks := strings.Split(key, "/")

    length := len(chunks)

    fmt.Fprintf(w, s.DataStore.Get(chunks[0]+"/"+chunks[1], chunks[length-1]))

    //log.Printf("[Debug]")
    //log.Printf(len(chunks))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/hello" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

    if r.Method != "GET" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }


    fmt.Fprintf(w, "Hello!")
}