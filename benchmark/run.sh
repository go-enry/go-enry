#!/bin/sh

benchmark/run-benchmark.sh && make benchmarks-slow && \
benchmark/parse.sh && benchmark/plot-histogram.sh
