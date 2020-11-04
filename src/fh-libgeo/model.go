package fh_libgeo

import (
	"time"
)

type ParseMetrics struct {
	Duration time.Duration
	InsertedOrUpdated int64
	Failed int64
	RowsProcessed int64
}

type GeolocationDbEntry struct {
	GeolocationDbEntryId int64 `json:"-" db:"geolocation_entry_id"` // do not send in json responses.
	InsertTime *time.Time `json:"-" db:"insert_time"` // do not send in json responses.
	IpAddress string `json:"ip_address" db:"ip_address"`
	CountryCode *string `json:"country_code" db:"country_code"`
	Country *string `json:"country" db:"country"`
	City *string `json:"city" db:"city"`
	Latitude *float64 `json:"latitude" db:"latitude"`
	Longitude *float64 `json:"longitude" db:"longitude"`
	MysteryValue *int64 `json:"mystery_value" db:"mystery_value"`
}

type GeolocationDbFailure struct {
	GeolocationDbEntryId int64 `json:"geolocation_entry_failure_id" db:"geolocation_entry_failure_id"`
	CsvFilename string `json:"csv_filename" db:"csv_filename"`
	CsvBytepos int64 `json:"csv_bytepos" db:"csv_bytepos"`
	CsvRow string `json:"csv_row" db:"csv_row"`
	FailureTime *time.Time `json:"geolocation_entry_failure_id" db:"failure_time"`
}
