package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/luchenyuxx/gophercise/url_shortener"
)

func main() {
	yamlFileName := flag.String("--yaml", "urlmapping.yaml", "the yaml file contains mapping")
	jsonFileName := flag.String("--json", "urlmapping.json", "the json file contains mapping")

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := url_shortener.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml, err := ioutil.ReadFile(*yamlFileName)
	if err != nil {
		panic(err)
	}
	yamlHandler, err := url_shortener.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	jsonFile, err := ioutil.ReadFile(*jsonFileName)
	if err != nil {
		panic(err)
	}
	jsonHandler, err := url_shortener.JSONHandler(jsonFile, yamlHandler)
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
