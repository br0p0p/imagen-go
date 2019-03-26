package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type ImageConfig struct {
	Width  uint64
	Height uint64
	Color  string "black"
}

func NewImageConfig() ImageConfig {
	image_config := ImageConfig{}
	image_config.Color = "black"
	return image_config
}

func Handler(w http.ResponseWriter, r *http.Request) {
	config := NewImageConfig()

	query_params := r.URL.Query()

	for k, v := range query_params {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}

	if width, err := strconv.ParseUint(query_params.Get("width"), 10, 64); err == nil {
		fmt.Printf("%T, %v\n", width, width)
		config.Width = width
	}

	if height, err := strconv.ParseUint(query_params.Get("height"), 10, 64); err == nil {
		fmt.Printf("%T, %v\n", height, height)
		config.Height = height
	}

	fmt.Printf("ImageConfig: %+v\n", config)

	fmt.Fprintf(w, "<h1>Hello from Go!")
}

func main() {
	http.HandleFunc("/", Handler)

	port := 8080

	// Handle static files
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println(fmt.Sprintf("Server listening at http://localhost:%d", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
