package ui

import (
	_ "embed"
	"log"
	"net/http"
)

//go:embed index.html
var html []byte

func ServeHTML(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write(html); err != nil {
		log.Fatal(err)
	}
}
