package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router := mux.NewRouter().StrictSlash(true)
	router.Path("/").Name("index").HandlerFunc(index)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func index(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("index.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	_, err = w.Write(file)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
