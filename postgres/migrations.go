package postgres

import (
	"dbtools/models"
	"fmt"
	"io/ioutil"
	"time"

	"dbtools/lib"
)

// CreateMigration --
func CreateMigration(app string, scname string, dbmodel string) error {
	tableprops := models.GetStruct(dbmodel, nil)

	pgdb.Init(app)
	res, err := pgdb.CreateMigrationUp(scname, dbmodel, tableprops)

	if err != nil {
		fmt.Println("Failed to create migration")
		return err
	}

	err = CreateMigrationFileUp(res, dbmodel)

	if err != nil {
		fmt.Println("Failed to create migration up file")
		return err
	}

	err = CreateMigrationFileDown(dbmodel)

	if err != nil {
		fmt.Println("Failed to create migration down file")
		return err
	}

	time.Sleep(1 * time.Second)

	return nil
}

// CreateMigrationFileUp --
func CreateMigrationFileUp(up string, name string) error {
	f := []byte(up)
	err := ioutil.WriteFile(fmt.Sprintf("database/migrations/%s_%s.up.sql", lib.CreateTimeStamp(), name), f, 0644)

	if err != nil {
		return err
	}

	return nil
}

// CreateMigrationFileDown --
func CreateMigrationFileDown(name string) error {
	f := []byte(fmt.Sprintf("DROP TABLE IF EXISTS %s;", name))
	err := ioutil.WriteFile(fmt.Sprintf("database/migrations/%s_%s.down.sql", lib.CreateTimeStamp(), name), f, 0644)

	if err != nil {
		return err
	}

	return nil
}

// CreateMigrationExternal --
func (p *Postgres) CreateMigrationExternal(filePath string, model interface{}) error {
	p.Connect()
	result, err := p.CreateMigrationUp(p.SchemaName, p.TableName, model)

	if err != nil {
		fmt.Println("Failed to create migration")
		return err
	}

	err = createMigrationFileUp(filePath, result, p.TableName)

	if err != nil {
		fmt.Println("Failed to create migration up file")
		return err
	}

	err = createMigrationFileDown(filePath, p.TableName)

	if err != nil {
		fmt.Println("Failed to create migration down file")
		return err
	}

	time.Sleep(1 * time.Second)
	p.Disconnect()
	return nil
}

func createMigrationFileUp(filePath string, up string, name string) error {
	byteFile := []byte(up)
	err := ioutil.WriteFile(fmt.Sprintf("%s/%s_%s.up.sql", filePath, lib.CreateTimeStamp(), name), byteFile, 0644)

	if err != nil {
		return err
	}

	return nil
}

func createMigrationFileDown(filePath string, name string) error {
	byteFile := []byte(fmt.Sprintf("DROP TABLE IF EXISTS %s", name))
	err := ioutil.WriteFile(fmt.Sprintf("%s/%s_%s.down.sql", filePath, lib.CreateTimeStamp(), name), byteFile, 0644)

	if err != nil {
		return err
	}

	return nil
}
