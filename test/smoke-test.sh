#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o xtrace

go run ./cmd/vib create profile test0
go run ./cmd/vib create profile test1
go run ./cmd/vib create profile test2
go run ./cmd/vib get profile
go run ./cmd/vib delete profile test{0,1,2}

echo Smoke tests ran successfully
