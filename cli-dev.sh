#!/bin/bash

# Read GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET from .env file
export $(cat .env | xargs)

START_DIR=$(pwd)

cd cli/

# Remove dist folder
rm -rf dist

go build -ldflags "-X main.version=dev -X TimeTrack/core/oauth.GOOGLE_CLIENT_ID=$GOOGLE_CLIENT_ID -X TimeTrack/core/oauth.GOOGLE_CLIENT_SECRET=$GOOGLE_CLIENT_SECRET" -o dist/TimeTrack

./dist/TimeTrack $@

cd $START_DIR