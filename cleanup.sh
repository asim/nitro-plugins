#!/bin/bash -e

for m in $(find ./ -name 'go.mod'); do
  d=$(dirname $m);
  if [ "$d" = "./micro/metrics" -o "$d" = "./micro/trace/awsxray" ]; then
    continue
  fi
  pushd $d;
  /bin/bash -c 'go get github.com/micro/go-micro/v2@v2.0.0 github.com/micro/micro/v2@v2.0.0 && go mod tidy' &
  popd;
done
