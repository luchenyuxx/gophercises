package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

type YAMLEntry struct {
	Path string
	Url  string
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func MapHandler(m map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for k, v := range m {
			log.Printf("%s:%s", k, v)
		}
		if val, ok := m[r.URL.Path]; ok {
			http.Redirect(w, r, val, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func buildMap(in []YAMLEntry) map[string]string {
	m := make(map[string]string)
	for _, entry := range in {
		m[entry.Path] = entry.Url
	}
	return m
}

func parseYAML(b []byte) ([]YAMLEntry, error) {
	var r []YAMLEntry
	e := yaml.Unmarshal(b, &r)
	return r, e
}
