#!/bin/bash -ex

mod=$(go list -m)
PKGS=""
for d in $(find * -name 'go.mod'); do
  PKGS=" $PKGS ${mod}/$(dirname $d)"
done

go test -race -v $PKGS || :
go test -v $PKGS
