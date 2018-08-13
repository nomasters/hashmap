#!/usr/bin/env bash
set -e

PROJ=hashmap-helper

## used to build the cli tool from source while killcord 
## was in private development

cd tools/$PROJ
go build
mv $PROJ $HOME/bin/$PROJ