DROP TABLE IF EXISTS geoentry CASCADE;

CREATE TABLE geoentry (
	geolocation_entry_id bigserial PRIMARY KEY NOT NULL,
	insert_time timestamp NOT NULL DEFAULT NOW(),
	ip_address text NOT NULL UNIQUE,
	country_code text,
	country text,
	city text,
	latitude double precision,
	longitude double precision,
	mystery_value bigint
);
ALTER TABLE geoentry OWNER TO fhgeo;

CREATE TABLE geoentry_csv_failures (
	geolocation_entry_failure_id bigserial PRIMARY KEY NOT NULL,
	csv_filename text NOT NULL,
	csv_bytepos integer NOT NULL,
	csv_row text NOT NULL,
	failure_time timestamp NOT NULL DEFAULT NOW()
);
ALTER TABLE geoentry_csv_failures OWNER TO fhgeo;

CREATE TABLE dbconfig (
	id	bigserial PRIMARY KEY NOT NULL,
	k	text NOT NULL UNIQUE,
	v	text NOT NULL DEFAULT '',
	update_time	timestamp NOT NULL DEFAULT NOW()
);
ALTER TABLE DbConfig OWNER TO fhgeo;

INSERT INTO dbconfig(k, v) VALUES ('schema_version', '1');
INSERT INTO dbconfig(k, v) VALUES ('data_version', '1');

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO fhgeo;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO fhgeo;

