#!/bin/sh

echo "PLAN: Adding ddvote database"
psql -U postgres -c 'CREATE DATABASE ddvote'

echo "PLAN: Creating tables"
psql -U postgres ddvote -f /dd-vote.sql
{ echo; echo "host all ddvote 0.0.0.0/0 trust"; } >> "$PGDATA"/pg_hba.conf
