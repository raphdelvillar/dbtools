package master

import (
	"dbtools/lib"
	"dbtools/postgres"
	"fmt"
	"reflect"
)

var (
	tables = map[string]interface{}{
		"test": []string{"sample-types"},
	}
)

// getStruct --
func getStruct(name string, data interface{}) interface{} {
	switch name {
	// test
	case "sample-types":
		return GetSampleType(data)
	}

	return nil
}

// CreateMaster --
func CreateMaster(app string, accountNumber string, scname string, pgdb postgres.Postgres) error {
	for _, table := range tables[scname].([]string) {
		err := seedPostgresWithMaster(app, fmt.Sprintf("%s-%s", accountNumber, scname), table, pgdb)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedPostgresWithMaster --
func seedPostgresWithMaster(app string, scname string, dbmodel string, pgdb postgres.Postgres) error {
	datas, tableprops, err := lib.ReadJSON(dbmodel, fmt.Sprintf("./master/%s/%s.json", app, dbmodel), getStruct)

	if err != nil {
		return err
	}

	err = pgdb.CreateSchema(scname)

	if err != nil {
		fmt.Printf("Failed to create schema %s", scname)
		return err
	}

	err = pgdb.CreateTable(scname, tableprops)

	if err != nil {
		fmt.Printf("Failed to create table %s", scname)
		return err
	}

	for _, data := range datas {
		delete(data.(map[string]interface{}), "id")
		structValue := getStruct(dbmodel, data)
		err := pgdb.Insert(structValue)

		if err != nil {
			fmt.Printf("Failed to insert %s", dbmodel)
			return err
		}

		createRelationship(scname, dbmodel, structValue, pgdb)
	}

	return nil
}

// createRelationship --
func createRelationship(scname string, dbmodel string, structValue interface{}, pgdb postgres.Postgres) {
	if structValue != nil {
		values := reflect.ValueOf(structValue).Elem()

		for i := 0; i < values.NumField(); i++ {
			value := values.Field(i)

			if value.Kind() == reflect.Struct {
				childStructValue := getStruct(fmt.Sprintf("%s", value.Type()), value.Interface())

				err := pgdb.CreateTable(scname, childStructValue)
				if err != nil {
					fmt.Println(err)
				}

				err = pgdb.Insert(childStructValue)
				if err != nil {
					fmt.Println(err)
				}

				createRelationship(scname, dbmodel, childStructValue, pgdb)
			}

			if value.Kind() == reflect.Array || value.Kind() == reflect.Slice {
				for i := 0; i < value.Len(); i++ {
					childStructValue := getStruct(fmt.Sprintf("%s", value.Index(i).Type()), value.Index(i).Interface())

					err := pgdb.CreateTable(scname, childStructValue)
					if err != nil {
						fmt.Println(err)
					}

					err = pgdb.Insert(childStructValue)
					if err != nil {
						fmt.Println(err)
					}

					createRelationship(scname, dbmodel, childStructValue, pgdb)
				}
			}
		}
	}
}
