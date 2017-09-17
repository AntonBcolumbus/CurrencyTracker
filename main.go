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
	router.Path("/getData").Name("index").HandlerFunc(getData)

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

func getData(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{
		"cols": [
			  {"id":"","label":"Topping","pattern":"","type":"string"},
			  {"id":"","label":"Slices","pattern":"","type":"number"}
			],
		"rows": [
			  {"c":[{"v":"Mushrooms","f":null},{"v":3,"f":null}]},
			  {"c":[{"v":"Onions","f":null},{"v":1,"f":null}]},
			  {"c":[{"v":"Olives","f":null},{"v":1,"f":null}]},
			  {"c":[{"v":"Zucchini","f":null},{"v":1,"f":null}]},
			  {"c":[{"v":"Pepperoni","f":null},{"v":2,"f":null}]}
			]
	  }`))
}
