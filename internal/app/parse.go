package app

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/vlatan/youtube-stats/web"
)

type fileInfo struct {
	bytes     []byte
	mediatype string
	Etag      string
}

// Create minified versions of the static files and cache them in memory.
func parseStaticFiles(m *minify.M, dir string) (map[string]fileInfo, error) {

	sf := make(map[string]fileInfo)

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
		b, err := fs.ReadFile(web.Files, path)
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
	if err := fs.WalkDir(web.Files, dir, walkDirFunc); err != nil {
		return nil, err
	}

	return sf, nil
}

// Custom function that minifies and parses the HTML templates
// as per the tdewolff/minify docs. Also inserts inline SVG image/map where needed.
func parseTemplates(
	m *minify.M,
	staticFiles map[string]fileInfo,
	filepaths ...string) (*template.Template, error) {

	m.AddFunc("text/html", html.Minify)

	var tmpl *template.Template
	for _, fp := range filepaths {

		b, err := fs.ReadFile(web.Files, fp)
		if err != nil {
			return nil, err
		}

		// inline the svg map if HTML id svgContainer present
		if svg, ok := staticFiles["map.svg"]; ok {
			htmlTag := []byte("id=\"svgContainer\">")
			inlineSVG := append(htmlTag, svg.bytes...)
			b = bytes.Replace(b, htmlTag, inlineSVG, 1)
		}

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
