package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/muesli/cache2go"
)

var cache *cache2go.CacheTable

const kb = 1024

func upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	uuid := uuid.New().String()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else if len(b) > 15*kb {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
	} else {
		cache.Add(uuid, 60*time.Minute, b)
		fmt.Fprint(w, uuid)
	}
}

func present(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res, err := cache.Value(ps.ByName("uuid"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		fmt.Fprintf(w, string(res.Data().([]byte)))
	}
}

func main() {
	cache = cache2go.Cache("cache")
	router := httprouter.New()
	router.POST("/", upload)
	router.GET("/:uuid", present)
	log.Fatal(http.ListenAndServe(":8080", router))
}
