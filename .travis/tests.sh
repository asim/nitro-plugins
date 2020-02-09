#!/bin/bash -ex

mod="github.com/micro/go-plugins"
PKGS=""
for d in $(find * -name 'go.mod'); do
  PKGS=" $PKGS ${mod}/$(dirname $d)/v2"
done

go test -race -v $PKGS || :
go test -v $PKGS
