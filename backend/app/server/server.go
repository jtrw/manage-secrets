package server

import (
 "fmt"
 //"log"
 "time"
 "net/http"

 log "github.com/go-pkgz/lgr"
 um "github.com/go-pkgz/rest"
 "github.com/pkg/errors"

 "github.com/didip/tollbooth/v6"
 "github.com/didip/tollbooth_chi"
 "github.com/go-chi/chi/v5"
 "github.com/go-chi/chi/v5/middleware"
)

type Server struct {
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
	router.Use(um.AppInfo("secrets", "Umputun", s.Version), um.Ping, um.SizeLimit(64*1024))
	router.Use(tollbooth_chi.LimitHandler(tollbooth.NewLimiter(10, nil)))

// 	router.Route("/api/v1", func(r chi.Router) {
// 		r.Use(Logger(log.Default()))
// 		r.Post("/message", s.saveMessageCtrl)
// 		r.Get("/message/{key}/{pin}", s.getMessageCtrl)
// 		r.Get("/params", s.getParamsCtrl)
// 	})
//
// 	router.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
// 		render.PlainText(w, r, "User-agent: *\nDisallow: /api/\nDisallow: /show/\n")
// 	})
//
// 	s.fileServer(router, "/", http.Dir(s.WebRoot))
	return router
}

// func Run() {
//     http.HandleFunc("/hello", helloHandler) // Update this line of code
//
//     fmt.Printf("Starting server at port 8080\n")
//     if err := http.ListenAndServe(":8080", nil); err != nil {
//         log.Fatal(err)
//     }
// }

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