#!/bin/bash -e

for m in $(find ./ -name 'go.mod'); do
  d=$(dirname $m);
  pushd $d;
  /bin/bash -c 'go get github.com/micro/cli/v2@master && go mod tidy' &
  popd;
done
