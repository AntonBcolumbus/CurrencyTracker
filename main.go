package main

import (
	"encoding/json"
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

// Dataset ...
type Dataset struct {
	Cols []Col `json:"cols,omitempty"`
	Rows []Row `json:"rows,omitempty"`
}

// Col ...
type Col struct {
	ID      string            `json:"id,omitempty"`
	Label   string            `json:"label,omitempty"`
	Type    string            `json:"type,omitempty"`
	Pattern string            `json:"pattern,omitempty"`
	P       map[string]string `json:"p,omitempty"`
	Role    string            `json:"role,omitempty"`
}

// Row ...
type Row struct {
	C []Cell `json:"c,omitempty"`
}

// Cell ...
type Cell struct {
	V interface{}       `json:"v,omitempty"`
	F string            `json:"f,omitempty"`
	P map[string]string `json:"p,omitempty"`
}

func getData(w http.ResponseWriter, r *http.Request) {

	data := Dataset{}
	data.Cols = []Col{
		Col{Label: "Day", Type: "number"},
		Col{Label: "Guardians of the Galaxy", Type: "number"},
		Col{Type: "string", P: map[string]string{"role": "annotation"}},
		Col{Type: "string", P: map[string]string{"role": "annotationText"}},
	}
	data.Rows = []Row{
		Row{
			C: []Cell{
				Cell{V: 1},
				Cell{V: 37.8},
				Cell{V: "37.8"},
				Cell{V: "37.8"},
			},
		},
		Row{
			C: []Cell{
				Cell{V: 2},
				Cell{V: 30.9},
				Cell{V: "30.9"},
				Cell{V: "30.9"},
			},
		},
		Row{
			C: []Cell{
				Cell{V: 3},
				Cell{V: 25.4},
				Cell{V: "25.4"},
				Cell{V: "25.4"},
			},
		},
		Row{
			C: []Cell{
				Cell{V: 4},
				Cell{V: 11.7},
				Cell{V: "11.7"},
				Cell{V: "11.7"},
			},
		},
		Row{
			C: []Cell{
				Cell{V: 5},
				Cell{V: 11.9},
				Cell{V: "11.9"},
				Cell{V: "11.9"},
			},
		},
		Row{
			C: []Cell{
				Cell{V: 6},
				Cell{V: 8.8},
				Cell{V: "8.8"},
				Cell{V: "8.8"},
			},
		},
	}

	d, _ := json.Marshal(data)
	w.Write(d)
}
