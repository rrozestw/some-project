#!/bin/bash
set -e

if [[ x"$1" == x"--force" ]] ; then
    sudo su postgres -c "psql -f db.sql"
    sudo su postgres -c "psql -d fhgeo -f tables.sql"
else
    echo "Recreates entire DB."
    echo "To confirm, execute with argument --force"
    exit 1
fi

