package app

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/vlatan/youtube-stats/internal/config"
	"github.com/vlatan/youtube-stats/internal/limiter"
	"github.com/vlatan/youtube-stats/internal/youtube"
	"golang.org/x/time/rate"
)

const staticFs = "static"
const tmplFile = "templates/index.html"

var validJS = regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$")

type App struct {
	youtube *youtube.Service
	static  map[string]fileInfo
	tmpl    *template.Template
	limiter *limiter.Service
	server  *http.Server
}

func New() (*App, error) {

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	ytsvc, err := youtube.New(context.Background(), cfg.ApiKey)
	if err != nil {
		return nil, err
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFuncRegexp(validJS, js.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)

	staticFiles, err := parseStaticFiles(m, staticFs)
	if err != nil {
		return nil, err
	}

	// Configure a custom HTML minifier
	htmlMinifier := &html.Minifier{
		KeepDocumentTags: true,                  // Prevent stripping <html>, <head>, and <body>
		KeepEndTags:      true,                  // Keep valid HTML structure
		TemplateDelims:   html.GoTemplateDelims, // Preserve context within golang template delimiters
	}

	// Use the custom HTML in a minifier function
	m.AddFunc("text/html", htmlMinifier.Minify)

	tmpl, err := parseTemplates(m, staticFiles, tmplFile)
	if err != nil {
		return nil, err
	}

	lim := limiter.New(rate.Every(time.Minute), 5, 5*time.Minute, 10*time.Minute)

	app := &App{
		youtube: ytsvc,
		static:  staticFiles,
		tmpl:    tmpl,
		limiter: lim,
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}

	app.server.Handler = app.Handler()
	return app, nil
}
