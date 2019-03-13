[![Build Status](https://travis-ci.org/mchirico/go_script.svg?branch=master)](https://travis-ci.org/mchirico/go_script)
[![Maintainability](https://api.codeclimate.com/v1/badges/9451bd1a6c801dd5eedb/maintainability)](https://codeclimate.com/github/mchirico/go_script/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/9451bd1a6c801dd5eedb/test_coverage)](https://codeclimate.com/github/mchirico/go_script/test_coverage)
[![codecov](https://codecov.io/gh/mchirico/go_script/branch/master/graph/badge.svg)](https://codecov.io/gh/mchirico/go_script)
[![Go Report Card](https://goreportcard.com/badge/github.com/mchirico/go_script)](https://goreportcard.com/report/github.com/mchirico/go_script)
# go_script
Trouble shooting bash script, managed by go.




## Libraries Necessary

```bash
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
go get github.com/prometheus/client_golang/prometheus/promhttp
go get -u github.com/gorilla/mux
go get gopkg.in/yaml.v2

```

## Build

```bash
mkdir -p scratch && cd scratch
git clone https://github.com/mchirico/go_script.git
cd go_script
go mod init 


go build ./cmd/script

# Yes, run this twice to create script file.
./script
./script

# Want to run tests?

go build ./...
go test ./...

```


# Build with vendor

```bash

mkdir -p scratch && cd scratch
git clone https://github.com/mchirico/go_script.git
cd go_script

go mod init
# Below will put all packages in a vendor folder
go mod vendor

export GO111MODULE=on

go test -v -mod=vendor ./...

# Don't forget the "." in "./cmd/script" below
go build -v -mod=vendor ./cmd/script



```


# Build for Linux on a Mac

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