package app

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type htmlData struct {
	StaticFiles map[string]fileInfo
}

// Handle GET request from the caller
func (a *App) Get(w http.ResponseWriter, r *http.Request) {

	data := htmlData{a.static}

	var buf bytes.Buffer
	if err := a.tmpl.ExecuteTemplate(&buf, "index.html", data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write to response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := buf.WriteTo(w); err != nil {
		log.Println(err)
	}
}

// Handle POST request by getting the video ID from the form,
// passing that ID to the YouTube API and returning a JSON
// response to the caller (writing to the http.ResponseWriter).
func (a *App) Post(w http.ResponseWriter, r *http.Request) {

	video, err := a.youtube.Video(r.Context(), r.FormValue("content"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// Encode data to JSON
	jsonData, err := json.Marshal(video)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write to response
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonData); err != nil {
		log.Println(err)
	}
}

// Handle minified static file from cache
func (a *App) Static(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	name := pathParts[len(pathParts)-1]
	fb, ok := a.static[name]

	// do not make the svg file accesable
	if !ok || strings.HasSuffix(name, ".svg") {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", fb.mediatype)
	w.Header().Set("Cache-Control", "max-age=31536000")
	w.Header().Set("Etag", fb.Etag)

	if _, err := w.Write(fb.bytes); err != nil {
		log.Println(err)
	}
}
