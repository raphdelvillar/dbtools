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
    dump
  else
    dump
  fi
}

function dump() {
  if [[ " ${APPS[*]} " == *"sample"* ]]
  then
    dumpdb "sample" HRM_DATABASES[@]
  fi
}

function dumpdb() {
    APP=${1}
    DATABASES=("${!2}")
    for database in "${DATABASES[@]}"; do
        echo "-----------------------------"
        echo -e "${CYAN}Dumping $database ...${NC}"
        COLLECTION="${APP^^}_${database^^}[@]"
        dumpcol "$ACCOUNT_NUMBER-${database}" $APP ${COLLECTION}
    done
}

function dumpcol() {
    COLAPP=${1}
    DIR=${2}
    COLLECTIONS=("${!3}")
    for collection in "${COLLECTIONS[@]}"; do
        echo "-----------------------------"
        echo -e "${CYAN}Dumping $collection ...${NC}"
        
        COLNAME=${collection}
        if [ "$4" != "" ]
        then
           COLNAME=${ACCOUNT_NUMBER}-${collection}
        fi

        echo ${COLAPP}
        echo ${COLNAME}

        COMMAND="mongoexport --host ${HOST} --db ${COLAPP} --collection=${COLNAME} -u ${USERNAME} -p ${PASSWORD} --authenticationDatabase admin -o dump/${DIR}/${collection}.json --jsonArray"
        eval $COMMAND
    done
}

run
echo "Done"