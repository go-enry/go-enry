#!/usr/bin/env bash
set -e

cd benchmarks/output
go run ../parser/main.go -outdir ../csv
cd ../csv
go run ../parser/main.go -distribution
