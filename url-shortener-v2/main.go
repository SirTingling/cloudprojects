package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var urlStore = make(map[string]string)

func shorten(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.Form.Get("url")
	shortUrl := fmt.Sprintf("%x", md5.Sum([]byte(url)))[:5] // take first 5 chars as key/s to map
	urlStore[shortUrl] = url
	fmt.Fprintf(w, "http://localhost:8080/%s\n", shortUrl)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	originalURL, ok := urlStore[vars["shortURL"]]
	if ok {
		http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/create", shorten).Methods("POST")
	r.HandleFunc("/{shortURL}", redirect).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
