package fh_libgeo

import (
	log "logex-gls"
	_ "database/sql"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

func getGeoEntryFromDb(db *sqlx.DB, ip string) (*GeolocationDbEntry, error) {
	geoentry := &GeolocationDbEntry{}
	err := db.Get(geoentry, "SELECT * FROM geoentry WHERE ip_address=$1", ip)
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	log.Infof("Got geoentry: %#v\n", geoentry)
	return geoentry, nil
}

func upsertNewGeoEntryToDb(db *sqlx.DB, e *GeolocationDbEntry) (int64, error) {
	var result int64
	err := db.QueryRow(`INSERT INTO geoentry(
	ip_address,
	country_code,
	country,
	city,
	latitude,
	longitude,
	mystery_value)
	 VALUES(
		$1, $2, $3, $4, $5, $6, $7
	) ON CONFLICT (ip_address)
		DO
			UPDATE
				SET country_code = $2,
					 country = $3,
					 city = $4,
					 latitude = $5,
					 longitude = $6,
					 mystery_value = $7,
					 insert_time = NOW()
	RETURNING geolocation_entry_id;`,
		e.IpAddress,
		e.CountryCode,
		e.Country,
		e.City,
		e.Latitude,
		e.Longitude,
		e.MysteryValue).Scan(&result)
	log.Debug(result)
	if err != nil {
		log.Error("FAILED adding new geoentry with error", err)
		return result, err
	}
	return result, nil
}

func insertNewGeoEntryToDb(db *sqlx.DB, e *GeolocationDbEntry) (int64, error) {
	var result int64
	err := db.QueryRow(`INSERT INTO geoentry(
	ip_address,
	country_code,
	country,
	city,
	latitude,
	longitude,
	mystery_value)
	 VALUES(
		$1, $2, $3, $4, $5, $6, $7
	) RETURNING geolocation_entry_id;`,
		e.IpAddress,
		e.CountryCode,
		e.Country,
		e.City,
		e.Latitude,
		e.Longitude,
		e.MysteryValue).Scan(&result)
	log.Debug(result)
	if err != nil {
		log.Error("FAILED adding new geoentry with error", err)
		return result, err
	}
	return result, nil
}

