#!/bin/sh

cd benchmark/output && go run ../parser/main.go -outdir ../csv && \
cd ../csv && go run ../parser/main.go -distribution

