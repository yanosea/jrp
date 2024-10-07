#!/usr/bin/env bash

# test packages and output a report file
go test -v -p 1 ../... -cover -coverprofile=./cover.out
# remove mock files from the coverage report
awk '!/(mock)\//' ./cover.out > temp_cover.out && mv temp_cover.out cover.out
# generate a html report for GitHub Pages
go tool cover -html=./cover.out -o ../docs/coverage.html
# remove the cover.out file
rm ./cover.out
