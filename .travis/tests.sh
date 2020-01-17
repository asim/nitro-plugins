#!/bin/bash

mod=$(go list -m)
PKGS=""
for d in $(find * -name 'go.mod'); do
  PKGS=" ${mod}/$(dirname $d)"
done

go test -race -v $PKGS || :
go test -v $PKGS
