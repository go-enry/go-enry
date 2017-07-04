#!/bin/sh

mkdir -p benchmark/output &&  go test -run NONE -bench=. -benchtime=120s -timeout=100h >benchmark/output/enry_total.bench && \
benchmark/linguist-total.rb 5 >benchmark/output/linguist_total.bench
