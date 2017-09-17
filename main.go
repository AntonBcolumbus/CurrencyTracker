package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

// TinkoffData ...
type TinkoffData struct {
	ResultCode string  `json:"resultCode,omitempty"`
	Payload    Payload `json:"payload,omitempty"`
	TrackingID string  `json:"trackingId,omitempty"`
}

// Payload ..
type Payload struct {
	LastUpdate LastUpdate `json:"lastUpdate,omitempty"`
	Rates      []Rate     `json:"rates,omitempty"`
}

// LastUpdate ...
type LastUpdate struct {
	Milliseconds uint64 `json:"milliseconds,omitempty"`
}

// Rate ...
type Rate struct {
	Category     string   `json:"category,omitempty"`
	FromCurrency Currency `json:"fromCurrency,omitempty"`
	ToCurrency   Currency `json:"toCurrency,omitempty"`
	Buy          float64  `json:"buy,omitempty"`
	Sell         float64  `json:"sell,omitempty"`
}

// Currency ...
type Currency struct {
	Name string `json:"name,omitempty"`
	Code int    `json:"code,omitempty"`
}

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.tinkoff.ru/v1/currency_rates", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	data := &TinkoffData{}
	json.Unmarshal(body, &data)

	url := ""
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("currency_charts").C("tinkoff")
	err = c.Insert(data)
	if err != nil {
		log.Fatal(err)
	}
}
