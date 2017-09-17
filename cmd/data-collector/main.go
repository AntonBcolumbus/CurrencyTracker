package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

func main() {

	loc, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(loc)

	if now.Weekday() == time.Sunday ||
		now.Weekday() == time.Saturday ||
		now.Hour() < 10 ||
		(now.Hour() > 18 && now.Minute() > 29) {
		log.Printf("skipping due to date/time: %s", now.Format("Mon Jan 2 15:04:05 MST 2006"))
		return
	}

	url := os.Getenv("MONGODB_URI")
	if url == "" {
		url = "mongodb://heroku_k99bcr9h:vo2e0n2drkk3do41t2q9lvh6av@ds141024.mlab.com:41024/heroku_k99bcr9h"
	}
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatalf("mongo connection failed: %s", err.Error())
	}
	defer session.Close()

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.tinkoff.ru/v1/currency_rates", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("tinkoff request failed: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("tinkoff request body reading failed: %s", err.Error())
	}

	data := &TinkoffData{}
	json.Unmarshal(body, &data)

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("heroku_k99bcr9h").C("tinkoff")
	err = c.Insert(data)
	if err != nil {
		log.Fatalf("mongo insert failed: %s", err.Error())
	}

	unixTime := now.UTC().Add(time.Hour*24*-3).UnixNano() / int64(time.Millisecond)
	err = c.Remove(bson.M{"payload.lastupdate.milliseconds": bson.M{"$lt": unixTime}})
	if err != nil && err != mgo.ErrNotFound {
		log.Fatalf("mongo remove failed: %s", err.Error())
	}
}
