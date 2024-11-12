package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"

	"github.com/joho/godotenv"
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
var home = template.Must(template.ParseFiles("web/templates/index.html"))

func main() {
	godotenv.Load()
	mux := http.NewServeMux()
	staticHandler := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", staticHandler))
	mux.HandleFunc("GET /{$}", getHandler)
	mux.HandleFunc("POST /{$}", postHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	err := home.Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var apiKey string = os.Getenv("YOUTUBE_API_KEY")
	if len(apiKey) == 0 {
		log.Fatal("Please set YOUTUBE_API_KEY environment variable.")
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
