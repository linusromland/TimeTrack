#!/bin/bash
export $(cat .env | xargs)

version=$1

if [ -z "$version" ]
then
	echo "Version is empty"
	exit 1
fi

START_DIR=$(pwd)


cd cli/

# Remove dist folder
rm -rf dist

go build -ldflags "-X main.version=$version -X TimeTrack/core/oauth.GOOGLE_CLIENT_ID=$GOOGLE_CLIENT_ID -X TimeTrack/core/oauth.GOOGLE_CLIENT_SECRET=$GOOGLE_CLIENT_SECRET" -o dist/TimeTrack

cd $START_DIR