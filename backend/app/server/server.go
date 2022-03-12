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

 "encoding/json"
)

type Server struct {
    DataStore      secret.Store
	PinSize        int
	MaxPinAttempts int
	WebRoot        string
	Version        string
}

//type JSON map[string]interface{}

func (s Server) Run() error {
    log.Printf("[INFO] Activate rest server")
    log.Printf("[INFO] Host: 127.0.0.1")
    log.Printf("[INFO] Port: 8080")

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
	    r.Post("/kv/*", s.saveValuesByKey)
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

func (s Server) saveValuesByKey(w http.ResponseWriter, r *http.Request) {
    b, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("[ERROR] %s", err)
    }
    value := string(b)
    log.Printf("[INFO] %s", value)

    uri := chi.URLParam(r, "*")
    keyStore, bucket := getKeyAndBucketByUrl(uri)

    dataJson := &secret.JSON{}
    dataType := "text"
    if isContentTypeJson(r) {
        errJsn := json.Unmarshal([]byte(value), dataJson)
        if errJsn != nil {
            log.Printf("ERROR Invalid json in Data");
            return
        }
        dataType = "json"
    }

    message := secret.Message {
        Key: keyStore,
        Bucket: bucket,
        Data: value,
        DataJson: *dataJson,
        Type: dataType,
    }

    s.DataStore.Save(&message)
    //s.DataStore.Set(bucket, keyStore, value)

    render.JSON(w, r, secret.JSON{"status": "ok"})
    return
}
func (s Server) getValuesByKey(w http.ResponseWriter, r *http.Request) {
    uri := chi.URLParam(r, "*")
    log.Printf("ContentType: %s", r.Header.Get("Content-Type"))
    onlyData := r.URL.Query().Get("onlyData")

    keyStore,bucket := getKeyAndBucketByUrl(uri)

    newMessage, _ := s.DataStore.Load(bucket, keyStore)

    render.Status(r, http.StatusCreated)
    if len(onlyData) > 0 {
        if isContentTypeJson(r) {
            render.JSON(w, r, newMessage.DataJson)
            return
        }
        render.JSON(w, r, newMessage.Data)
    } else {
        render.JSON(w, r, newMessage)
    }


    //render.JSON(w, r, JSON{"key": newMessage.Key, "Data": newMessage.Data})
    //render.JSON(w, r, JSON{"Data": json.newMessage.Data})
}


func getKeyAndBucketByUrl(uri string) (string, string) {
    chunks := strings.Split(uri, "/")

    length := len(chunks)
    keyStore := chunks[length-1]
    bucket := strings.Join(chunks[:length-1], "/")

    return keyStore, bucket
}

func isContentTypeJson(r *http.Request) bool {
    return r.Header.Get("Content-Type") == strings.ToLower("application/json")
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