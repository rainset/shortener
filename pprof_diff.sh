#!/bin/bash
go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof
