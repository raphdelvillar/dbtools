#!/bin/bash

# DEFAULT PATHS
BASE_DIR=$(pwd)

source "$BASE_DIR/mongo/.env"

if [ "$1" != "" ];
then
  APPS=($1)
else
  APPS=()
fi

function run() {
  if [ ${#APPS[@]} -eq 0 ];
  then
    APPS=("sample")
    go
  else
    go
  fi
}

function go() {
  mkdir models
  if [[ " ${APPS[*]} " == *"sample"* ]]
  then
    mkdir models/sample
    creatego "sample"
  fi
}

function creatego() {
    path="./schema/$1/*"
    for f in $path; do
        filename=$(basename "${f%.*}")
        echo "-----------------------------"
        echo -e "${CYAN}Generating Go Struct $filename ...${NC}"
        mkdir "./models/$1/$filename"
        COMMAND="quicktype --package $filename -s schema $f -o ./models/$1/$filename/$filename.go"
        eval $COMMAND
    done
}

run
echo "Done"