package lib

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler that read url mapping from map
func MapHandler(m map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for k, v := range m {
			fmt.Printf("%s:%s", k, v)
		}
		if val, ok := m[r.URL.Path]; ok {
			http.Redirect(w, r, val, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// JSONHandler that read url mapping from yaml format
func JSONHandler(j []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(j)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

func parseJSON(j []byte) ([]urlMap, error) {
	var r []urlMap
	err := json.Unmarshal(j, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// YAMLHandler that read url mapping from yaml format
func YAMLHandler(y []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(y)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

type urlMap struct {
	Path string
	URL  string
}

func buildMap(in []urlMap) map[string]string {
	m := make(map[string]string)
	for _, entry := range in {
		m[entry.Path] = entry.URL
	}
	return m
}

func parseYAML(b []byte) ([]urlMap, error) {
	var r []urlMap
	e := yaml.Unmarshal(b, &r)
	if e != nil {
		return nil, e
	}
	return r, nil
}
