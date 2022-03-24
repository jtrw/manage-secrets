package server

import (
 "fmt"
 "io"
 //"lgr"
 "time"
 "net/http"
 "strings"
 "crypto/rand"
 "encoding/hex"
 "log"
 lgr "github.com/go-pkgz/lgr"
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
const ENV_TOKEN_KEY = "APP_JTRW_SECRET_TOKEN"

type Server struct {
    DataStore      secret.Store
    Host           string
    Port           string
	PinSize        int
	MaxPinAttempts int
	WebRoot        string
	Version        string
}

//type JSON map[string]interface{}

func (s Server) Run() error {
	if err := http.ListenAndServe(s.Host+":"+s.Port, s.routes()); err != http.ErrServerClosed {
		return errors.Wrap(err, "server failed")
	}
	return nil
}

func (s Server) routes() chi.Router {
	router := chi.NewRouter()

    token := GenerateSecureToken(20)

    fmt.Printf("Please add this token to .env file. Property %s \n", ENV_TOKEN_KEY)
    fmt.Printf("Token: %s \n", token)
    s.saveToken(token)

    router.Use(middleware.Logger)
    router.Use(s.AuthMiddleware)
	router.Use(middleware.RequestID, middleware.RealIP, um.Recoverer(lgr.Default()))
	router.Use(middleware.Throttle(1000), middleware.Timeout(60*time.Second))
	router.Use(um.AppInfo("secrets", "jtrw", s.Version), um.Ping, um.SizeLimit(64*1024))
	router.Use(tollbooth_chi.LimitHandler(tollbooth.NewLimiter(10, nil)))

	router.Route("/api/v1", func(r chi.Router) {
	    r.Get("/kv/*", s.getValuesByKey)
	    r.Post("/kv/*", s.saveValuesByKey)
	   // r.Post("/token/", s.saveToken)
		//r.Use(Logger(lgr.Default()))
		//r.Get("/message/{key}/{pin}", s.getMessageCtrl)
	})

    lgr.Printf("[INFO] Activate rest server")
    lgr.Printf("[INFO] Host: %s", s.Host)
    lgr.Printf("[INFO] Port: %s", s.Port)

	return router
}

func (s Server) AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
       token := s.getToken()
       accessToken := r.Header.Get("Access-Token")
       if token != accessToken {
            http.Error(rw, http.StatusText(403), 403)
            return
       }
       next.ServeHTTP(rw, r)
	})
}

func GenerateSecureToken(length int) string {
    b := make([]byte, length)
    if _, err := rand.Read(b); err != nil {
        return ""
    }
    return hex.EncodeToString(b)
}

func (s Server) saveToken(token string) {
    message := secret.Message {
        Key: "token",
        Bucket: secret.TOKEN_KEY,
        Data: token,
    }

    s.DataStore.Save(&message)
}

func (s Server) getToken() string {
     message, err := s.DataStore.Load(secret.TOKEN_KEY, "token")
     if err != nil {
        log.Fatal("Token Not Found!")
     }
     return message.Data
}

func (s Server) saveValuesByKey(w http.ResponseWriter, r *http.Request) {
    b, err := io.ReadAll(r.Body)
    if err != nil {
        lgr.Printf("[ERROR] %s", err)
    }
    value := string(b)
    lgr.Printf("[INFO] %s", value)

    uri := chi.URLParam(r, "*")
    keyStore, bucket := getKeyAndBucketByUrl(uri)

    dataJson := &secret.JSON{}
    dataType := "text"
    if isContentTypeJson(r) {
        errJsn := json.Unmarshal([]byte(value), dataJson)
        if errJsn != nil {
            lgr.Printf("ERROR Invalid json in Data");
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
    render.Status(r, http.StatusCreated)
    render.JSON(w, r, secret.JSON{"status": "ok"})
    return
}
func (s Server) getValuesByKey(w http.ResponseWriter, r *http.Request) {
    uri := chi.URLParam(r, "*")
    lgr.Printf("ContentType: %s", r.Header.Get("Content-Type"))
    onlyData := r.URL.Query().Get("onlyData")

    keyStore,bucket := getKeyAndBucketByUrl(uri)

    newMessage, _ := s.DataStore.Load(bucket, keyStore)

    render.Status(r, http.StatusOK)
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