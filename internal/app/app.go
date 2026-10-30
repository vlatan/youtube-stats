package app

import (
	"context"
	"html/template"
	"net/http"
	"regexp"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/vlatan/youtube-stats/internal/youtube"
)

const staticFs = "static"
const tmplFile = "templates/index.html"

var validJS = regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$")

type App struct {
	youtube *youtube.Service
	static  map[string]fileInfo
	tmpl    *template.Template
	server  *http.Server
}

func New(ctx context.Context, apiKey string) (*App, error) {

	ytsvc, err := youtube.New(ctx, apiKey)
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

	return &App{
		youtube: ytsvc,
		static:  staticFiles,
		tmpl:    tmpl,
	}, nil
}
