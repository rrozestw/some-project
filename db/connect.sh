#!/bin/bash -x
CONN_STR='postgresql://fhgeo:1234@localhost/fhgeo?sslmode=disable'
pgcli "$CONN_STR"
