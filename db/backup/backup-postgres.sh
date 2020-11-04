#!/bin/bash -x

CURRENT_DATE=$(date -u +"%Y-%m-%dT%H-%M-%SZ")
sudo -u postgres pg_dump fhgeo | gzip - > ./fhgeo_$CURRENT_DATE.sql.gz
# RESTORE: 
# 1. sudo su postgres -c "psql -f db.sql"
# 2. zcat fhgeo_*.sql.gz | sudo -u postgres psql fhgeo
