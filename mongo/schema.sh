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
    schema
  else
    schema
  fi
}

function schema() {
  mkdir schema
  if [[ " ${APPS[*]} " == *"sample"* ]]
  then
    mkdir schema/sample
    createschema "sample"
  fi
}

function createschema() {
    path="./dump/$1/*"
    for f in $path; do
        filename=$(basename "${f%.*}")
        echo "-----------------------------"
        echo -e "${CYAN}Generating Schema $filename ...${NC}"
        COMMAND="quicktype $f -l schema -o ./schema/$1/$filename.json"
        eval $COMMAND
    done
}

run
echo "Done"