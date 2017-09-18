package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"

	"github.com/AntonBcolumbus/CurrencyTracker/datamodels"
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

	data := getMongoData()

	dataToSend := Dataset{}
	dataToSend.Cols = []Col{
		Col{Label: "Day", Type: "datetime"},
		Col{Label: "EUR - RUB", Type: "number"},
		Col{Type: "string", P: map[string]string{"role": "annotation"}},
		Col{Label: "USD - RUB", Type: "number"},
		Col{Type: "string", P: map[string]string{"role": "annotation"}},
	}
	dataToSend.Rows = make([]Row, 0, 0)

	for _, d := range data {
		eurC := Cell{}
		usdC := Cell{}
		var eurFound, usdFound bool = false, false
		for _, r := range d.Payload.Rates {
			if r.Category == "SMETransferBelow10" && r.FromCurrency.Code == 978 && r.ToCurrency.Code == 643 {
				eurC.V = r.Buy
				eurFound = true
				continue
			}
			if r.Category == "SMETransferBelow10" && r.FromCurrency.Code == 840 && r.ToCurrency.Code == 643 {
				usdC.V = r.Buy
				usdFound = true
				continue
			}
			if usdFound && eurFound {
				break
			}
		}
		row := Row{}
		t := time.Unix(0, d.Payload.LastUpdate.Milliseconds*int64(time.Millisecond)).UTC()
		loc, _ := time.LoadLocation("Europe/Moscow")
		t = t.In(loc)
		row.C = []Cell{
			Cell{
				V: fmt.Sprintf("Date(%d,%d,%d,%d,%d)", t.Year(), t.Month()-1, t.Day(), t.Hour(), t.Minute()),
			},
			eurC,
			Cell{
				V: fmt.Sprintf("%v", eurC.V),
			},
			usdC,
			Cell{
				V: fmt.Sprintf("%v", usdC.V),
			},
		}
		dataToSend.Rows = append(dataToSend.Rows, row)
	}

	d, _ := json.Marshal(dataToSend)
	w.Write(d)
}

func getMongoData() []*datamodels.TinkoffData {
	url := os.Getenv("MONGODB_URI")
	if url == "" {
		url = "mongodb://heroku_k99bcr9h:vo2e0n2drkk3do41t2q9lvh6av@ds141024.mlab.com:41024/heroku_k99bcr9h"
	}
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatalf("mongo connection failed: %s", err.Error())
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("heroku_k99bcr9h").C("tinkoff")
	var result []*datamodels.TinkoffData
	loc, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(loc)

	millisFrom := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, loc).UnixNano() / int64(time.Millisecond)
	millisTo := time.Date(now.Year(), now.Month(), now.Day(), 18, 30, 0, 0, loc).UnixNano() / int64(time.Millisecond)
	err = c.Find(bson.M{
		"$and": []bson.M{
			bson.M{"payload.lastupdate.milliseconds": bson.M{"$gte": millisFrom}},
			bson.M{"payload.lastupdate.milliseconds": bson.M{"$lte": millisTo}},
		}}).All(&result)
	if err != nil {
		log.Fatalf("mongo find failed: %s", err.Error())
	}
	return result
}
