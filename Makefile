#MONGO SPECIFIC COMMANDS
mongo-dump: 
	./mongo/dump.sh $(apps)
mongo-schema: 
	./mongo/schema.sh $(apps)
mongo-go: 
	./mongo/go.sh $(apps)
mongo-create: mongo-dump mongo-schema mongo-go

#CREATE MIGRATIONS
#This will create a timestamp_schema.sql file under database/migrations folder
create-migration:
	go run *.go mg $(app) $(accountNumber) $(schema)

#CREATE SEEDS
#This will create a timestamp_schema.sql file under database/seeds folder
create-seed:
	go run *.go sd $(app) $(accountNumber) $(schema)

#MONGO TO POSTGRES
#This will create a dump folder with mongo data then populate the db according to the models
mongo-postgres: mongo-dump
	go run *.go ps $(app) $(accountNumber) $(schema)

#POSTGRES TO MONGO
#This will get data from postgres then save that data to mongo
postgres-mongo:
	go run *.go ms $(app) $(accountNumber) $(schema)

#MONGO TO POSTGRES MASTER
#This will populate the db with master data
mongo-postgres-master:
	go run *.go cm $(app) $(accountNumber) $(schema)

#MIGRATION
# create_migration:
# 	migrate create -ext sql -dir database/migrations $(name)

# create_seed:
# 	migrate create -ext sql -dir database/seeds $(name)

up_migration:
	@./migrate.sh up

down_migration:
	@./migrate.sh down

version_migration:
	@./migrate.sh version

drop_migration:
	@./migrate.sh drop