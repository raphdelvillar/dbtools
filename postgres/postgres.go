package postgres

import (
	"fmt"
	"strings"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/jinzhu/inflection"
)

var pgdb Postgres

// Postgres --
type Postgres struct {
	Auth       pg.Options
	DB         *pg.DB
	SchemaName string
	TableName  string
}

// Init --
func (p *Postgres) Init(app string) {
	p.Auth = pg.Options{
		Addr:     "",
		User:     "",
		Password: "",
		Database: app,
	}

	p.Connect()
}

// Connect --
func (p *Postgres) Connect() {
	db := pg.Connect(&p.Auth)
	p.DB = db
}

// Disconnect --
func (p *Postgres) Disconnect() {
	p.DB.Close()
}

// CreateDatabase --
func (p *Postgres) CreateDatabase(dbname string) error {
	command := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbname)
	_, err := p.DB.Exec(command)

	if err != nil {
		return err
	}

	return nil
}

// CreateMigrationUp --
func (p *Postgres) CreateMigrationUp(scname string, tablename string, model interface{}) (string, error) {
	err := p.CreateSchema(scname)

	if err != nil {
		return "", err
	}

	err = p.CreateTable(scname, model)

	if err != nil {
		fmt.Println(err)
	}

	query := fmt.Sprintf(`with c as (
		SELECT table_name, ordinal_position,
		 column_name|| ' ' || data_type col
		, row_number() over (partition by table_name order by ordinal_position asc) rn
		, count(*) over (partition by table_name) cnt
		FROM information_schema.columns
		WHERE table_schema = '%s' AND
		table_name = '%s'
		order by table_name, ordinal_position
		)
		select case when rn = 1 then 'CREATE TABLE IF NOT EXISTS "' || concat('%s.', table_name) || '" (' else '' end
		 || col
		 || case when rn < cnt then ',' else '); ' end AS query
		from c
		order by table_name, rn asc;`, scname, tablename, scname)

	var result pg.Strings

	_, err = p.DB.Query(&result, query)

	if err != nil {
		return "", err
	}

	err = p.DropTable(scname, model)

	if err != nil {
		return "", err
	}

	return strings.Join(result, "\n"), nil
}

// CreateSeed --
func (p *Postgres) CreateSeed(scname string, tablename string, model interface{}) (string, error) {
	err := p.CreateSchema(scname)

	if err != nil {
		return "", err
	}

	err = p.CreateTable(scname, model)

	if err != nil {
		fmt.Println(err)
	}

	query := fmt.Sprintf(`with c as (
		SELECT table_name, ordinal_position, 
		 column_name col
		, row_number() over (partition by table_name order by ordinal_position asc) rn
		, count(*) over (partition by table_name) cnt
		FROM information_schema.columns
		WHERE table_schema = '%s' AND 
		table_name = '%s'
		order by table_name, ordinal_position
		)
		select case when rn = 1 then 'INSERT INTO "' || concat('%s.', table_name) || '" (' else '' end
		 || col 
		 || case when rn < cnt then ',' else ')' end AS result
		from c
		order by table_name, rn asc;`, scname, tablename, scname)

	var result pg.Strings

	_, err = p.DB.Query(&result, query)

	if err != nil {
		return "", err
	}

	err = p.DropTable(scname, model)

	if err != nil {
		return "", err
	}

	return strings.Join(result, "\n"), nil
}

// CreateSchema --
func (p *Postgres) CreateSchema(scname string) error {
	schemaName := fmt.Sprintf("\"%s\"", scname)
	command := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", schemaName)
	_, err := p.DB.Exec(command)

	if err != nil {
		return err
	}

	return nil
}

// SetTableNameInflector --
func (p *Postgres) setTableNameInflector(scname string) {
	orm.SetTableNameInflector(func(s string) string {
		return fmt.Sprintf("%s.", scname) + inflection.Plural(s)
	})
}

// CreateTable --
func (p *Postgres) CreateTable(scname string, model interface{}) error {
	p.setTableNameInflector(scname)
	err := p.DB.CreateTable(model, &orm.CreateTableOptions{
		Temp: false,
	})
	if err != nil {
		return err
	}
	return nil
}

// DropTable --
func (p *Postgres) DropTable(scname string, model interface{}) error {
	p.setTableNameInflector(scname)
	err := p.DB.DropTable(model, &orm.DropTableOptions{
		IfExists: true,
	})
	if err != nil {
		return err
	}
	return nil
}

// Insert --
func (p *Postgres) Insert(data interface{}) error {
	err := p.DB.Insert(data)
	if err != nil {
		return err
	}
	return nil
}

// CountRow --
func (p *Postgres) CountRow(scname string, data interface{}) (int, error) {
	p.setTableNameInflector(scname)
	count, err := p.DB.Model(data).Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// SelectOne --
func (p *Postgres) SelectOne(scname string, data interface{}, offset int) (interface{}, error) {
	p.setTableNameInflector(scname)
	err := p.DB.Model(data).Offset(offset).Limit(1).Select()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// SelectMany --
func (p *Postgres) SelectMany(scname string, data interface{}) {
	p.setTableNameInflector(scname)
	err := p.DB.Model(data).Select()
	if err != nil {
		fmt.Println(err)
	}
}
