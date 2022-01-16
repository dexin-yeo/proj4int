#!/bin/bash

PROJDIR=$(dirname $(pwd))
BIN=$PROJDIR/bin

if [ -f "$BIN/app" ]; then
	rm $BIN/app
fi
