#!/bin/bash
source .env

# DEFAULT PATHS
BASE_DIR=$(pwd)

migrate -source file://${BASE_DIR}/database/migrations/ -database "$DATABASE"?sslmode=disable $1 $2