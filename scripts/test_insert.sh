#!/bin/bash

PROJDIR=$(dirname $(pwd))
APP=$PROJDIR/bin/app
CSV=$PROJDIR/data/table-definition.csv

$APP -i "f9cd42b93c434beba00f97bcbdd20dbe" $CSV

$APP -i "f9cd42b93c434beba00f97bcbdd20dbe" "my_table" "tmp1" "18" "" "NO" "integer"
$APP -i "f9cd42b93c434beba00f97bcbdd20dbe" "my_table" "tmp2" "19" "" "NO" "character varying" "66"