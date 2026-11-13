// Go-Web-App-Skeleton
/*
# Build

    $ go build

# Run

    $ ./go-webskeleton --port 7000 &

*/
package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

var port = flag.Int("port", 8080, "Set another port: --port 9000")
var tmpl *template.Template

//go:embed static
//go:embed templates/*.html
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

	fs := http.FileServer(http.FS(staticFS))

	http.Handle("/static/", fs)
	http.HandleFunc("/", mainHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {

	path, err := exec.LookPath("uptime")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	b, err := exec.Command(path, "-p").Output()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	s, err := exec.Command(path, "-s").Output()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := struct {
		Port   int
		Uptime string
		Since  string
	}{
		Port:   *port,
		Uptime: string(b[2:]),
		Since:  string(s),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		return
	}
}
