#!/bin/bash -ex

mod=$(go list -m | sed 's|/v2||g')
PKGS=""
for d in $(find * -name 'go.mod'); do
  pushd $(dirname $d)
  go mod download
  popd
done
