package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type story []string

type storyArc struct {
	Title   string   `json:"title"`
	Story   story    `json:"story"`
	Options []option `json:"options"`
}

type cyoaHandler struct {
	stories map[string]storyArc
}

func (o cyoaHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
	default:
		w.Write([]byte("hello"))
	}
}

func main() {
	gopherFile := "gopher.json"
	jsonFile, err := ioutil.ReadFile(gopherFile)
	if err != nil {
		panic(err)
	}
	var cyoa map[string]storyArc
	err = json.Unmarshal(jsonFile, &cyoa)
	if err != nil {
		panic(err)
	}
	handler := cyoaHandler{cyoa}
	http.ListenAndServe(":8080", handler)
}
