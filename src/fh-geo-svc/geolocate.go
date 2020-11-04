package main

import (
	"os"
	"time"
	"net/http"
	"encoding/json"
	"github.com/gocraft/web"
	log "logex-gls"
	libgeo "fh-libgeo"
	_ "database/sql"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"fmt"
)

func getLibGeo() (libgeo.LibGeo, error) {
	db, err := sqlx.Open("postgres", "user="+os.Getenv("POSTGRES_USER")+" dbname="+
		os.Getenv("POSTGRES_DB")+" password="+os.Getenv("POSTGRES_PASS")+" sslmode="+os.Getenv("POSTGRES_SSLMODE")+" host="+os.Getenv("POSTGRES_HOST")+" port="+os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return nil, err
	}
	lib, err := libgeo.New(db)
	return lib, err
}

// panics on failure - recovery not implemented. Would need to keep the last row/byte/line processed and start from there.
func csv_loader() {
	log.Info("CSV loader...")
	libgeo, err := getLibGeo()
	if err != nil {
		log.Error(err)
		panic(err)
	}
	for {
		result, err := libgeo.InsertToDbFromCsvFile("data/data_dump.csv")
		if err != nil {
			log.Error(err)
			panic(err)
		}
		fmt.Println("Duration:", result.Duration)
		fmt.Println("InsertedOrUpdated:", result.InsertedOrUpdated)
		fmt.Println("Failed:", result.Failed)
		time.Sleep(1 * time.Minute)
	}
}

func (c *Context) GeolocateEndpoint(rw web.ResponseWriter, req *web.Request) {
	req_ip := req.PathParams["ip"]
	log.Info("geolocate ", req_ip)
	libgeo, err := getLibGeo()
	if err != nil {
		log.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	lookup_result, err := libgeo.Lookup(req_ip)
	if err != nil {
		log.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	lookup_json_bytes, _ := json.Marshal(lookup_result)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(lookup_json_bytes)
}

