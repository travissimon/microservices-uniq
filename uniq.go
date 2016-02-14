package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/travissimon/remnant/client"
)

const (
	SERVICE_PATH = "/v1/uniq/"
)

func uniqHandler(w http.ResponseWriter, r *http.Request, cl client.Client) {
	cl.LogDebug("request: " + r.URL.Path)
	s := string(r.URL.Path[len(SERVICE_PATH):])
	charMap := make(map[rune]bool)
	for _, c := range s {
		charMap[c] = true
	}

	var buffer bytes.Buffer
	for k, _ := range charMap {
		buffer.WriteRune(k)
	}
	var allChars = buffer.String()

	cl.LogInfo(fmt.Sprintf("Unique chars in %s: %s", s, allChars))

	// count those chars
	destUrl := "http://localhost:8001/v1/strlen/" + allChars
	resp, err := cl.Get(destUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling %s: %s\n", destUrl, err.Error())
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling %s: %s\n", destUrl, err.Error())
		return
	}

	w.Write(body)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("healthz %s\n", time.Now().Local())
	fmt.Fprintf(w, "OK")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "No really a thing")
}

func main() {
	var port = flag.String("port", "8080", "Define which TCP port to bind to")
	flag.Parse()

	remnantUrl := "http://localhost:7777/"
	http.HandleFunc(SERVICE_PATH, client.GetInstrumentedHandler(remnantUrl, uniqHandler))
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/", indexHandler)

	fmt.Printf("Starting uniq server on port %s\n", *port)
	http.ListenAndServe(":"+*port, nil)
}
