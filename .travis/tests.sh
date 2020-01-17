#!/bin/bash

mod=$(go list -m)

for d in $(find * -name 'go.mod'); do
  go test -race -v ${mod}/$(dirname $d) || :
  go test -v ${mod}/$(dirname $d)
done
