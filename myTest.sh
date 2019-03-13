#!/bin/bash
cd ~/go_script
. setpath

if [  -f go.mod ]; then
    rm go.mod
fi

if [ -d vendor ]; then
   rm -rf vendor
fi


export GO111MODULE=on
go mod init
go mod vendor


go fmt ./...
go test -race -v -mod=vendor -coverprofile=c.out  ./...  && echo -e "\n\n\n ✅ SUCCESS \n\n"
