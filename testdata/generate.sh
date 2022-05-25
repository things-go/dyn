#!/bin/bash

PROTOS=$(find ./encoding -type f -name '*.proto')

for PROTO in $PROTOS; do
  echo $PROTO
  protoc \
    -I . \
    -I $(dirname $PROTO) \
    -I ../third_party \
    --go_out=. \
    --go_opt=paths=source_relative \
    $PROTO
done