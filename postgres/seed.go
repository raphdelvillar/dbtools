package postgres

import (
	"dbtools/lib"
	"dbtools/models"
	"fmt"
	"reflect"
)

// SeedPostGresWithMongo --
func SeedPostGresWithMongo(app string, scname string, dbmodel string) error {
	datas, tableprops, err := lib.ReadJSON(dbmodel, fmt.Sprintf("./dump/%s/%s.json", app, dbmodel), models.GetStruct)

	if err != nil {
		fmt.Println(err)
	}

	pgdb.Init(app)
	err = pgdb.CreateSchema(scname)

	if err != nil {
		fmt.Printf("Failed to create schema %s", scname)
		fmt.Println(err)
	}

	err = pgdb.CreateTable(scname, tableprops)

	if err != nil {
		fmt.Printf("Failed to create table %s", scname)
		fmt.Println(err)
	}

	for _, data := range datas {
		delete(data.(map[string]interface{}), "id")
		CreatePostGresData(scname, dbmodel, data)
	}

	return nil
}

// CreatePostGresData --
func CreatePostGresData(scname string, dbmodel string, data interface{}) {
	structValue := models.GetStruct(dbmodel, data)
	err := pgdb.Insert(structValue)

	if err != nil {
		fmt.Printf("Failed to insert %s", dbmodel)
		fmt.Println(err)
	}

	createRelationship(scname, dbmodel, structValue)
}

// createRelationship --
func createRelationship(scname string, dbmodel string, structValue interface{}) {
	if structValue != nil {
		values := reflect.ValueOf(structValue).Elem()

		for i := 0; i < values.NumField(); i++ {
			value := values.Field(i)

			if value.Kind() == reflect.Struct {
				childStructValue := models.GetStruct(fmt.Sprintf("%s", value.Type()), value.Interface())

				err := pgdb.CreateTable(scname, childStructValue)
				if err != nil {
					fmt.Println(err)
				}

				err = pgdb.Insert(childStructValue)
				if err != nil {
					fmt.Println(err)
				}

				createRelationship(scname, dbmodel, childStructValue)
			}

			if value.Kind() == reflect.Array || value.Kind() == reflect.Slice {
				for i := 0; i < value.Len(); i++ {
					childStructValue := models.GetStruct(fmt.Sprintf("%s", value.Index(i).Type()), value.Index(i).Interface())

					err := pgdb.CreateTable(scname, childStructValue)
					if err != nil {
						fmt.Println(err)
					}

					err = pgdb.Insert(childStructValue)
					if err != nil {
						fmt.Println(err)
					}

					createRelationship(scname, dbmodel, childStructValue)
				}
			}
		}
	}
}
