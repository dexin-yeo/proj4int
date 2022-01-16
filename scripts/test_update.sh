#!/bin/bash

PROJDIR=$(dirname $(pwd))
APP=$PROJDIR/bin/app
CSV=$PROJDIR/data/sampleCSV.csv

$APP -u $CSV