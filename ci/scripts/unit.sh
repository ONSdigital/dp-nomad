#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-nomad
  make test
popd