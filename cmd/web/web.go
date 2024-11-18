package main

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
	common "github.com/vlatan/youtube-stats/internal"
)

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

var validID = regexp.MustCompile("^([-a-zA-Z0-9_]{11})$")
var validJS = regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$")
var home = template.Must(template.ParseFiles("web/templates/index.html"))

func main() {

	staticDir := "web/static"

	err := minifyStaticFiles(staticDir)
	if err != nil {
		log.Fatal(err)
	}

	staticHandler := http.FileServer(http.Dir(staticDir))
	staticHandler = http.StripPrefix("/static/", staticHandler)

	godotenv.Load()
	mux := http.NewServeMux()

	mux.Handle("GET /static/", staticHandler)
	mux.HandleFunc("GET /{$}", getHandler)
	mux.HandleFunc("POST /{$}", postHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

// Create minified versions of the static files on disk within the same
// static directory with ".min" inserted before their file extension.
func minifyStaticFiles(root string) error {

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFuncRegexp(validJS, js.Minify)

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
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// split the file path on dot
		pathParts := strings.Split(path, ".")

		// set media type (just css or js)
		mediatype := "text/css"
		if pathParts[1] == "js" {
			mediatype = "application/javascript"
		}

		// minify the content
		b, err = m.Bytes(mediatype, b)
		if err != nil {
			return err
		}

		// insert "min" into the path and write to disk
		minPath := pathParts[0] + ".min." + pathParts[1]
		err = os.WriteFile(minPath, b, 0644)
		if err != nil {
			return err
		}

		return nil
	}

	return filepath.WalkDir(root, walkDirFunc)

}

// Handle GET request from the caller
// by loading a HTML template to the response.
func getHandler(w http.ResponseWriter, r *http.Request) {
	err := home.Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
}

// Handle POST request by getting the video ID from the form,
// passing that ID to the YouTube API and returning a JSON
// response to the caller (writing to the http.ResponseWriter).
func postHandler(w http.ResponseWriter, r *http.Request) {
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
