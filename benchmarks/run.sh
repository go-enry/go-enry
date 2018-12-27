#!/usr/bin/env bash
set -e

benchmarks/run-benchmarks.sh
make benchmarks-slow
benchmarks/parse.sh
benchmarks/plot-histogram.gp
