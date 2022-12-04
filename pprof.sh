#!/bin/bash
export PPROF_TMPDIR="$HOME/Sites/go/shortener/profiles"
go tool pprof -http=":9090" -seconds=30  http://localhost:8080/debug/pprof/heap

