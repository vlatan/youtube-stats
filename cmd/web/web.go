package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"

	"github.com/joho/godotenv"
	common "github.com/vlatan/youtube-stats/internal"
)

type Video struct {
	Id               string   `json:"id,omitempty"`
	Title            string   `json:"title,omitempty"`
	PrivacyStatus    string   `json:"privacyStatus,omitempty"`
	AgeRestriced     bool     `json:"ageRestriced,omitempty"`
	Embeddable       bool     `json:"embeddable,omitempty"`
	RegionRestricted []string `json:"regionRestricted,omitempty"`
	DefaultLanguage  string   `json:"defaultLanguage,omitempty"`
	LiveBroadcast    string   `json:"liveBroadcast,omitempty"`
	Duration         string   `json:"duration,omitempty"`
}

var validID = regexp.MustCompile("^([-a-zA-Z0-9_]+)$")

// var templates = template.Must(template.ParseGlob(filepath.Join("web", "templates", "*.html")))
var home = template.Must(template.ParseFiles("web/templates/index.html"))
var content = template.Must(template.ParseFiles("web/templates/index.html", "web/templates/content.html"))

func main() {
	godotenv.Load()
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", makeHandler(rootHandler))
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func rootHandler(w http.ResponseWriter, r *http.Request, apiKey string) {
	// if r.URL.Path != "/" {
	// 	http.Redirect(w, r, "/", http.StatusFound)
	// }

	id := r.FormValue("id")
	if id == "" {
		loadTemplate(w, home, nil)
		return
	}

	if validID.FindStringSubmatch(id) == nil {
		log.Println("Not a valid ID")
		loadTemplate(w, home, nil)
		return
	}

	response, err := common.GetVideo(apiKey, id)
	if err != nil {
		log.Println(err)
		loadTemplate(w, home, nil)
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

	loadTemplate(w, content, video)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var apiKey string = os.Getenv("YOUTUBE_API_KEY")
		if len(apiKey) == 0 {
			log.Fatal("Please set YOUTUBE_API_KEY environment variable.")
		}
		fn(w, r, apiKey)
	}
}

func loadTemplate(w http.ResponseWriter, tmpl *template.Template, data any) {
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
}
