package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	common "github.com/vlatan/youtube-stats/internal"
)

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

// Handle GET request from the caller
func getHandler(w http.ResponseWriter, r *http.Request) {

	token, err := getCSRFToken(w, r)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}

	// Include the CSRF token in data to be passed to the template
	data := htmlData{
		CSRFToken:     token,
		CSRFFieldName: csrfFieldName,
		StaticFiles:   staticFiles,
	}

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

	// Get the the HttpOnly cookie
	if !validateCSRFToken(r) {
		http.Error(w, "CSRF token validation failed", http.StatusForbidden)
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

	url := r.FormValue("content")
	videoID, err := common.ExtractYouTubeID(url)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	response, err := common.GetVideo(apiKey, videoID)
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
