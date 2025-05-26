package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
	resources "github.com/vlatan/youtube-stats"
)

var (
	m           = minify.New()
	validID     = regexp.MustCompile("^([-a-zA-Z0-9_]{11})$")
	validJS     = regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$")
	staticFiles = parseStaticFiles(m, "web/static")
	templates   = template.Must(parseTemplates(m, "web/templates/index.html"))
)

type fileInfo struct {
	bytes     []byte
	mediatype string
	Etag      string
}

type cachedStaticFiles map[string]fileInfo

func getIntEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	port, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Invalid value for %s: %s. Using default port %d.", key, value, fallback)
		return fallback
	}

	return port
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
