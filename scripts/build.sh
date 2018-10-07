#!/usr/bin/env bash
set -e

PROJ=hashmap
## used to build the cli tool from source while in private development
cd $PROJ
go build
mkdir -p $HOME/bin/
mv $PROJ $HOME/bin/$PROJ
