package postgres

import (
	"dbtools/models"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"dbtools/lib"
)

// CreateSeed --
func CreateSeed(app string, scname string, dbmodel string) error {
	datas, tableprops, err := lib.ReadJSON(dbmodel, fmt.Sprintf("./dump/%s/%s.json", app, dbmodel), models.GetStruct)

	pgdb.Init(app)
	res, err := pgdb.CreateSeed(scname, dbmodel, tableprops)

	if err != nil {
		return err
	}

	values, err := CreateValues(datas, dbmodel, res)

	if err != nil {
		return err
	}

	err = CreateSeedFile(values, dbmodel)

	if err != nil {
		return err
	}

	return nil
}

// CreateValues --
func CreateValues(datas []interface{}, dbmodel string, insert string) ([]string, error) {
	var qValues []string
	for _, data := range datas {
		structValue := models.GetStruct(dbmodel, data)
		values := reflect.ValueOf(structValue).Elem()

		var qValue []string
		for i := 0; i < values.NumField(); i++ {
			value := values.Field(i)

			switch value.Kind() {
			case reflect.String:
				qValue = append(qValue, fmt.Sprintf("'%s'", value.String()))
			case reflect.Int:
				qValue = append(qValue, fmt.Sprintf("%v", strconv.FormatInt(value.Int(), 10)))
			case reflect.Int32:
				qValue = append(qValue, fmt.Sprintf("%v", strconv.FormatInt(value.Int(), 10)))
			case reflect.Int64:
				qValue = append(qValue, fmt.Sprintf("%v", strconv.FormatInt(value.Int(), 10)))
			default:
				qValue = append(qValue, fmt.Sprintf("'%v'", value.Interface()))
			}
		}

		query := fmt.Sprintf("%s VALUES (%s)", insert, strings.Join(qValue, ","))

		qValues = append(qValues, query)
	}

	return qValues, nil
}

// CreateSeedFile --
func CreateSeedFile(seed []string, name string) error {
	f := []byte(strings.Join(seed, "\n\n"))
	err := ioutil.WriteFile(fmt.Sprintf("database/seeds/%s_%s.sql", lib.CreateTimeStamp(), name), f, 0644)

	if err != nil {
		return err
	}

	return nil
}
