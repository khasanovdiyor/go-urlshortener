package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/khasanovdiyor/go-urlshortener/urlshortener"
)

func main() {
	yamlFile := flag.CommandLine.String("yaml", "redirects.yaml", "a path to yaml file (default is 'redirects.yaml')")
	jsonFile := flag.CommandLine.String("json", "redirects.json", "a path to json file (default is redirects.json)")

	flag.Parse()

	if flag.Lookup("h") != nil {
		flag.Usage()
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshortener.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlFileContents, err := os.ReadFile(*yamlFile)
	if err != nil {
		panic(err)
	}
	yamlHandler, err := urlshortener.YAMLHandler(yamlFileContents, mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler using the yamlHandler as the
	// fallback
	jsonFileContents, err := os.ReadFile(*jsonFile)
	if err != nil {
		panic(err)
	}
	jsonHandler, err := urlshortener.JSONHandler(jsonFileContents, yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
