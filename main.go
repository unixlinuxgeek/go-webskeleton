// Go-Web-App-Skeleton

package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var port = flag.Int("port", 8080, "Set another port: --port 9000")
var tmpl *template.Template

//go:embed static/*
//go:embed templates/*
var staticFS embed.FS

func main() {
	flag.Parse()

	fmt.Printf("port: %v\n", *port)
	var err error

	tmpl, err = template.ParseFS(staticFS, "templates/main.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "tmpl: %v\n", err)
		os.Exit(1)
	}

	fs := http.FS(staticFS)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(fs)))
	http.HandleFunc("/", mainHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, nil)
	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		return
	}
}
