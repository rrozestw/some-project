package fh_libgeo

import (
	log "logex-gls"
	"errors"
	"os"
	"io"
	"strconv"
	"time"
	_ "database/sql"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"bufio"
	"encoding/csv"
	"github.com/davecgh/go-spew/spew"
)

type LibGeo interface {
	Lookup(ip string) (GeolocationDbEntry, error)
	InsertToDbFromCsvFile(filepath string) (ParseMetrics, error)
}

type LibGeoPostgres struct {
	Db *sqlx.DB
}

func New(db *sqlx.DB) (LibGeo, error) {
	ret := &LibGeoPostgres{}
	if db == nil {
		return nil, errors.New("Database connector may not be nil.")
	}
	ret.Db = db
	return ret, nil
}

func (libgeo *LibGeoPostgres) Lookup(ip string) (GeolocationDbEntry, error) {
	entry, err := getGeoEntryFromDb(libgeo.Db, ip)
	if entry == nil {
		return GeolocationDbEntry{}, err
	} else {
		return *entry, err
	}
}

func (libgeo *LibGeoPostgres) InsertToDbFromCsvFile(filepath string) (ParseMetrics, error) {
	// TODO: split files on \n (read till byte around rough partition boundry and seek till \n is found; split) and process them in parallel for _, p := range parts { go process_part(p) }
	//  with above, how to handle upserts? Save row id in db, upsert only if row id > current row id.
	csvFile, err := os.Open(filepath)
	if err != nil {
		log.Println("Unable to open file:", filepath, err)
		return ParseMetrics{}, nil
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = ','
	//reader.Comment = '#'
	reader.FieldsPerRecord = 7
	// reader.LazyQuotes = 
	// reader.TrimLeadingSpace
	ret := ParseMetrics{ InsertedOrUpdated: 0, RowsProcessed: 0, Failed: 0 }
	row_id := 0
	time_begin := time.Now()
	for {
		row_id += 1
		ret.RowsProcessed += 1
		if row_id == 1 {
			continue // skipping header. This may not be the case all the time.
		}
		if row_id % 100 == 0 {
			elapsed := time.Now().Sub(time_begin)
			log.Warn("Parsing row", row_id, "elapsed time:", elapsed, "avg performance[rows/s]: ", float64(float64(row_id)/elapsed.Seconds()))
		}
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error("Error parsing file:", filepath, err)
			return ParseMetrics{}, nil
		}
		log.Debug("Parsing row:", row_id, row)
		entry := GeolocationDbEntry{}
		if len(row[0]) > 0 {
			entry.IpAddress = row[0]
		} else {
			log.Error("makes no sense to add an record without ip address.", row_id, row)
			ret.Failed += 1
			continue
		}
		// TODO: trim strings?
		if len(row[1]) == 2 {
			// assuming https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
			entry.CountryCode = &row[1]
		} else {
			log.Warn("Invalid country code:", row[1], row_id, row)
			// TODO: Try to find country code from country via lookup table;
		}
		if len(row[2]) > 0 {
			entry.Country = &row[2]
		}
		if len(row[3]) > 0 {
			entry.City = &row[3]
		}
		latitude, err_lat := strconv.ParseFloat(row[4], 64)
		longitude, err_long := strconv.ParseFloat(row[5], 64)
		if err_lat == nil && err_long == nil {
			// TODO: limit range to -180. ... 180.
			entry.Latitude = &latitude
			entry.Longitude = &longitude
		} else {
			log.Warn("Missing at least one coord, skipping both.", row_id, row)
		}
		mystery_value, err := strconv.ParseInt(row[6], 10, 64)
		if err == nil {
			entry.MysteryValue = &mystery_value
		}
		log.Info("row_id: ", row_id, "ENTRY:", spew.Sdump(entry))
		res_id, err := upsertNewGeoEntryToDb(libgeo.Db, &entry) 
		if err != nil {
			log.Fatal(err) // this should not fail, except when we have database connectivity errors. TODO retries for db failures?
			panic(err)
		}
		log.Debug("inserted id:", res_id, "; csv row id: ", row_id)
		ret.InsertedOrUpdated += 1
	}
	time_end := time.Now()
	ret.Duration = time_end.Sub(time_begin)
	log.Info("Parser metrics:", spew.Sdump(ret))
	return ret, nil
}

