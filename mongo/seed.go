package mongo

import (
	"dbtools/models"
	"fmt"
	"strings"
)

var (
	dbmigration = map[string]interface{}{
		"sample": "sample",
	}
)

// SeedMongoWithPostGres --
func SeedMongoWithPostGres(app string, scname string, dbmodel string) error {
	tableProps := models.GetStruct(dbmodel, nil)
	pgdb.Init(app)

	rowCount, err := pgdb.CountRow(scname, tableProps)

	if err != nil {
		return err
	}

	mdb.Init()

	schemaName := strings.Split(scname, "-")
	mdbSchemaName := fmt.Sprintf("%s-%s-%s", schemaName[0], schemaName[1], dbmigration[schemaName[2]])

	mdb.Database = mdbSchemaName
	mdb.Collection = dbmodel

	for i := 0; i < rowCount; i++ {
		data, err := pgdb.SelectOne(scname, tableProps, i)
		if err != nil {
			return err
		}
		mdb.Insert(data)
	}

	return nil
}
