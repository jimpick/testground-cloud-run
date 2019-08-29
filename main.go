package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Home page request.")
	body, err := ioutil.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	trace := r.Header.Get("X-Cloud-Trace-Context")
	log.Print("Test request ", trace)

	var lines []string

	line := fmt.Sprintf("Test")
	lines = append(lines, line)
	io.WriteString(w, strings.Join(lines, "\n"))
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/test", testHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8099"
	}

	log.Printf("Web server started on port: %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

