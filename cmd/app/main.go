package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/joho/godotenv"
	"github.com/tdewolff/minify/v2"
	"golang.org/x/time/rate"
)

type cachedStaticFiles map[string]fileInfo
type htmlData struct {
	CSRFToken     string
	CSRFFieldName string
	StaticFiles   cachedStaticFiles
}

type Video struct {
	Id               string   `json:"id"`
	Title            string   `json:"title"`
	PrivacyStatus    string   `json:"privacyStatus"`
	AgeRestriced     bool     `json:"ageRestriced"`
	Embeddable       bool     `json:"embeddable"`
	RegionRestricted []string `json:"regionRestricted"`
	DefaultLanguage  string   `json:"defaultLanguage"`
	LiveBroadcast    string   `json:"liveBroadcast"`
	Duration         string   `json:"duration"`
}

type fileInfo struct {
	bytes     []byte
	mediatype string
	Etag      string
}

const csrfCookieName = "csrf_token"
const csrfFieldName = "csrf_token"

var (
	m           = minify.New()
	validID     = regexp.MustCompile("^([-a-zA-Z0-9_]{11})$")
	validJS     = regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$")
	staticFiles = parseStaticFiles(m, "web/static")
	templates   = template.Must(parseTemplates(m, "web/templates/index.html"))
	limiter     = newIPRateLimiter(rate.Every(time.Minute), 5, 5*time.Minute, 10*time.Minute)
)

func main() {
	godotenv.Load()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /static/", staticHandler)
	mux.HandleFunc("GET /{$}", getHandler)

	postHandler := applyMiddlewares(postHandler, limitMiddleware)
	mux.HandleFunc("POST /{$}", postHandler)

	port := getIntEnv("PORT", 8080)
	fmt.Printf("Server running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
