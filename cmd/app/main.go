package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
	resources "github.com/vlatan/youtube-stats"
	common "github.com/vlatan/youtube-stats/internal"
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

// Handle minified static file from cache
func staticHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	name := pathParts[len(pathParts)-1]
	fb, ok := staticFiles[name]

	// do not make the svg file accesable
	if !ok || strings.HasSuffix(name, ".svg") {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", fb.mediatype)
	w.Header().Set("Cache-Control", "max-age=31536000")
	w.Header().Set("Etag", fb.Etag)

	if _, err := w.Write(fb.bytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Handle GET request from the caller
func getHandler(w http.ResponseWriter, r *http.Request) {

	csrfToken, err := generateCSRFToken()
	if err != nil {
		log.Println(err)
	}

	// Include the CSRF token in data to be passed to the template
	data := htmlData{
		CSRFToken:     csrfToken,
		CSRFFieldName: csrfFieldName,
		StaticFiles:   staticFiles,
	}

	// Set the token in a HttpOnly cookie
	setCSRFCookie(w, csrfCookieName, csrfToken)

	err = templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
}

// Handle POST request by getting the video ID from the form,
// passing that ID to the YouTube API and returning a JSON
// response to the caller (writing to the http.ResponseWriter).
func postHandler(w http.ResponseWriter, r *http.Request) {

	// Get the token from the HttpOnly cookie
	cookieToken, err := getCSRFCookie(r, csrfCookieName)
	if err != nil || cookieToken == "" {
		http.Error(w, "CSRF cookie not found or invalid", http.StatusForbidden)
		return
	}

	// Get the token from the form submission
	formToken := r.FormValue(csrfFieldName)

	// Compare the tokens
	if formToken != cookieToken {
		http.Error(w, "CSRF token mismatch!", http.StatusForbidden)
		return
	}

	// This is going to be a JSON response
	w.Header().Set("Content-Type", "application/json")

	var apiKey string = os.Getenv("YOUTUBE_API_KEY")
	if len(apiKey) == 0 {
		log.Println("Please set YOUTUBE_API_KEY environment variable.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := r.FormValue("id")
	if validID.FindStringSubmatch(id) == nil {
		log.Println("Not a valid video ID:", id)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	response, err := common.GetVideo(apiKey, id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	video := Video{
		Id:               response.Id,
		Title:            response.Snippet.Title,
		PrivacyStatus:    response.Status.PrivacyStatus,
		AgeRestriced:     common.AgeRestriction(response),
		Embeddable:       response.Status.Embeddable,
		RegionRestricted: common.RegionRestriction(response),
		DefaultLanguage:  response.Snippet.DefaultLanguage,
		LiveBroadcast:    response.Snippet.LiveBroadcastContent,
		Duration:         response.ContentDetails.Duration,
	}

	// write JSON to response
	if json.NewEncoder(w).Encode(video) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Create minified versions of the static files and cache them in memory.
func parseStaticFiles(m *minify.M, dir string) cachedStaticFiles {

	sf := cachedStaticFiles{}

	m.AddFunc("text/css", css.Minify)
	m.AddFuncRegexp(validJS, js.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)

	// function used to process each file/dir in the root, including the root
	walkDirFunc := func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// skip directories
		if info.IsDir() {
			return nil
		}

		// skip minified files
		if strings.Contains(info.Name(), ".min.") {
			return nil
		}

		// read the file
		b, err := fs.ReadFile(resources.Files, path)
		if err != nil {
			return err
		}

		// split the file path on dot
		pathParts := strings.Split(path, ".")

		// set media type
		mediatype := "text/css"
		switch pathParts[1] {
		case "js":
			mediatype = "application/javascript"
		case "svg":
			mediatype = "image/svg+xml"
		}

		// minify the content
		b, err = m.Bytes(mediatype, b)
		if err != nil {
			return err
		}

		// create Etag as a hexadecimal md5 hash of the file content
		etag := fmt.Sprintf("%x", md5.Sum(b))

		// save all the data in the struct
		sf[info.Name()] = fileInfo{b, mediatype, etag}

		return nil
	}

	// walk the directory and process each file
	if err := fs.WalkDir(resources.Files, dir, walkDirFunc); err != nil {
		log.Println(err)
	}

	return sf
}

// Custom function that minifies and parses the HTML templates
// as per the tdewolff/minify docs. Also inserts inline SVG image/map where needed.
func parseTemplates(m *minify.M, filepaths ...string) (*template.Template, error) {

	m.AddFunc("text/html", html.Minify)

	var tmpl *template.Template
	for _, fp := range filepaths {

		b, err := fs.ReadFile(resources.Files, fp)
		if err != nil {
			return nil, err
		}

		// inline the svg map if HTML id svgContainer present
		svg := staticFiles["map.svg"].bytes
		htmlTag := []byte("id=\"svgContainer\">")
		svg = append(htmlTag, svg...)
		b = bytes.Replace(b, htmlTag, svg, 1)

		name := filepath.Base(fp)
		if tmpl == nil {
			tmpl = template.New(name)
		} else {
			tmpl = tmpl.New(name)
		}

		mb, err := m.Bytes("text/html", b)
		if err != nil {
			return nil, err
		}
		tmpl.Parse(string(mb))
	}
	return tmpl, nil
}
