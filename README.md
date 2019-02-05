[![Build Status](https://travis-ci.org/mchirico/go_script.svg?branch=master)](https://travis-ci.org/mchirico/go_script)
[![Maintainability](https://api.codeclimate.com/v1/badges/9451bd1a6c801dd5eedb/maintainability)](https://codeclimate.com/github/mchirico/go_script/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/9451bd1a6c801dd5eedb/test_coverage)](https://codeclimate.com/github/mchirico/go_script/test_coverage)
[![codecov](https://codecov.io/gh/mchirico/go_script/branch/master/graph/badge.svg)](https://codecov.io/gh/mchirico/go_script)
# go_script
Go script to run troubleshooting commands without creating large logs


## Libraries Necessary

```bash
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
go get github.com/prometheus/client_golang/prometheus/promhttp
go get -u github.com/gorilla/mux

```

## Build

```bash
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
go build -o scriptLinux github.com/mchirico/go_script/cmd/script

```

## Table of Values

```bash
$GOOS     $GOARCH
darwin    386      – 32 bit MacOSX
darwin    amd64    – 64 bit MacOSX
freebsd   386
freebsd   amd64
linux     386      – 32 bit Linux
linux     amd64    – 64 bit Linux
linux     arm      – RISC Linux
netbsd    386
netbsd    amd64
openbsd   386
openbsd   amd64
plan9     386
windows   386      – 32 bit Windows
windows   amd64    – 64 bit Windows
```