package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
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

	cmd := exec.Command("./smlbench")
	cmd.Dir = "/home/testground/testground/plans/smlbench/cmd"

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Start()
	if err != nil {
		http.Error(w, "500 Exception", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	var errStdout, errStderr error

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
		wg.Done()
	}()

	_, errStderr = io.Copy(stderr, stderrIn)
	wg.Wait()

	err = cmd.Wait()
	if err != nil || errStdout != nil || errStderr != nil {
		http.Error(w, "500 Exception", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	output := strings.Split(string(stdoutBuf.Bytes()), "\n")
	for _, line := range output {
		lines = append(lines, line)
	}
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

