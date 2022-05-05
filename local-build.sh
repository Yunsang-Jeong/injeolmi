#!/bin/bash

GIT_ROOT=$(git rev-parse --show-toplevel)

TF_DIR=$GIT_ROOT/infra

cd $GIT_ROOT
GOOS=linux GOARCH=amd64 go build -o $TF_DIR/main main.go

cd $TF_DIR
zip app.zip main -x "*.DS_Store"
rm $TF_DIR/main