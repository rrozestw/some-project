DROP DATABASE IF EXISTS fhgeo;
DROP USER IF EXISTS fhgeo;

CREATE USER fhgeo WITH PASSWORD '1234';

CREATE DATABASE fhgeo WITH OWNER fhgeo;
GRANT ALL PRIVILEGES 
ON DATABASE fhgeo TO fhgeo;
