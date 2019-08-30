# What dbtools can do ?
- create postgres sql migration files from go structs
- create postgres sql seeder files from mongo dumps
- seed postgres db with mongo dumps 

# Library Required
go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate

# CREATE MIGRATIONS

This will create a timestamp_schema.*.sql file under database/migrations folder

make create-migration $(app) $(schema)

# RUN MIGRATIONS

make up_migration<br />
make down_migration<br />
make version_migration<br />
make drop_migration<br />

# CREATE SEEDS

This will create a timestamp_schema.sql file under database/seeds folder

make create-seed $(app) $(schema)

# MONGO TO POSTGRES

This will create a dump folder with mongo data then populate the db according to the models

make mongo-postgres $(app) $(schema)
