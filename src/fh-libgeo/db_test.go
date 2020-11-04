package fh_libgeo

import (
	"log"
	"testing"
	util "fh-common"
	_ "database/sql"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"os"
	"github.com/davecgh/go-spew/spew"
)

func setup() {
	util.CheckOrSetEnv("POSTGRES_DB", "fhgeo")
	util.CheckOrSetEnv("POSTGRES_USER", "fhgeo")
	util.CheckOrSetEnvPassword("POSTGRES_PASS", "1234")
	util.CheckOrSetEnv("POSTGRES_HOST", "localhost")
	util.CheckOrSetEnv("POSTGRES_SSLMODE", "disable")
	util.CheckOrSetEnv("POSTGRES_PORT", "5432")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func TestUpsert(t *testing.T) {
	setup()
	db, err := sqlx.Open("postgres", "user="+os.Getenv("POSTGRES_USER")+" dbname="+os.Getenv("POSTGRES_DB")+" password="+os.Getenv("POSTGRES_PASS")+" sslmode="+os.Getenv("POSTGRES_SSLMODE")+" host="+os.Getenv("POSTGRES_HOST")+" port="+os.Getenv("POSTGRES_PORT"))
	check(err)
	
	cc := "PL"
	country := "Poland"
	city := "Gdynia"
	entry := GeolocationDbEntry{
		IpAddress: "127.0.0.1",
		CountryCode: &cc,
		Country: &country,
		City: &city,
	}
	
	row_id, err := upsertNewGeoEntryToDb(db, &entry)
	check(err)

	ret_e, err := getGeoEntryFromDb(db, "127.0.0.1")
	check(err)
	log.Println("row_id: ", row_id, "ENTRY:", spew.Sdump(ret_e))
	
	cc = "NL"
	country = "Netherlands"
	city = "Amsterdam"
	entry = GeolocationDbEntry{
		IpAddress: "127.0.0.1",
		CountryCode: &cc,
		Country: &country,
		City: &city,
	}
	
	row_id2, err := upsertNewGeoEntryToDb(db, &entry)
	check(err)
	if row_id != row_id2 {
		panic("row IDs are different")
	}
	
	ret_e2, err := getGeoEntryFromDb(db, "127.0.0.1")
	check(err)
	log.Println("row_id: ", row_id2, "ENTRY:", spew.Sdump(ret_e2))
	
	if *ret_e2.CountryCode == *ret_e.CountryCode {
		panic("cc not updated;")
	}
	if *ret_e2.CountryCode != "NL" {
		panic("cc after upsert - wrong value.")
	}
	//panic("showstdout")
}
