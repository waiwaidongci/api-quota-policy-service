#!/usr/bin/env sh
set -eu
gofmt -w .
go test ./...
