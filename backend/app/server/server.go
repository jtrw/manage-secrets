package server

import (
 "fmt"
 "io"
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
 "github.com/go-chi/render"

 secret "manager-secrets/backend/app/store"
)

type Server struct {
    DataStore      secret.Store
	PinSize        int
	MaxPinAttempts int
	WebRoot        string
	Version        string
}

type JSON map[string]interface{}

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

	router.Route("/api/v1", func(r chi.Router) {
	    r.Get("/kv/*", s.getValuesByKey)
	    r.Post("/kv/*", s.setValuesByKey)
		//r.Use(Logger(log.Default()))
		//r.Get("/message/{key}/{pin}", s.getMessageCtrl)
	})
//
// 	router.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
// 		render.PlainText(w, r, "User-agent: *\nDisallow: /api/\nDisallow: /show/\n")
// 	})
//
// 	s.fileServer(router, "/", http.Dir(s.WebRoot))
	return router
}

func (s Server) setValuesByKey(w http.ResponseWriter, r *http.Request) {
    b, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("[ERROR] %s", err)
    }
    value := string(b)
    log.Printf("[INFO] %s", value)

    uri := chi.URLParam(r, "*")
    keyStore, bucket := getKeyAndBucketByUrl(uri)

//     message := secret.Message {
//         Key: keyStore,
//         Bucket: bucket,
//         Data: value,
//     }
//
//     s.DataStore.Save(&message)
    s.DataStore.Set(bucket, keyStore, value)

    render.JSON(w, r, JSON{"status": "ok"})
    return
}
func (s Server) getValuesByKey(w http.ResponseWriter, r *http.Request) {
    uri := chi.URLParam(r, "*")

    keyStore,bucket := getKeyAndBucketByUrl(uri)

    log.Printf("[INFO] %s", bucket)
    log.Printf("[INFO] %s", keyStore)
//      message := secret.Message {
//             Key: keyStore,
//             Bucket: bucket,
//         }
    //newMessage, _ := s.DataStore.Load(&message)
    //fmt.Fprintf(w, newMessage.Data)
    fmt.Fprintf(w, s.DataStore.Get(bucket, keyStore))
}


func getKeyAndBucketByUrl(uri string) (string, string) {
    chunks := strings.Split(uri, "/")

    length := len(chunks)
    keyStore := chunks[length-1]
    bucket := strings.Join(chunks[:length-1], "/")

    return keyStore, bucket
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