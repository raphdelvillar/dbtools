package main

import (
	"dbtools/master"
	"dbtools/mongo"
	"dbtools/postgres"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"
)

var (
	pgtables = map[string]interface{}{
		"sample": []string{"sample-types", "sample-model"},
	}

	mstables = map[string]interface{}{
		"sample": []string{"sample"},
	}
)

type fn func(string, string, string) error

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			// To run "go run *.go mg"
			Name:    "create-migrations",
			Aliases: []string{"mg"},
			Usage:   "command -- mg [app] [schema]",
			Action: func(c *cli.Context) error {
				if os.Args[2] != "" && os.Args[3] != "" {
					err := execute(os.Args[2], "", os.Args[3], postgres.CreateMigration, pgtables)

					if err != nil {
						panic(err)
					}
				}

				return nil
			},
		},
		{
			// To run "go run *.go sd"
			Name:    "create-seeds",
			Aliases: []string{"sd"},
			Usage:   "command -- sd [app] [schema]",
			Action: func(c *cli.Context) error {
				if os.Args[2] != "" && os.Args[3] != "" {
					err := execute(os.Args[2], "", os.Args[3], postgres.CreateSeed, pgtables)

					if err != nil {
						panic(err)
					}
				}

				return nil
			},
		},
		{
			// To run "go run *.go ps"
			Name:    "postgres-seed",
			Aliases: []string{"ps"},
			Usage:   "command -- ps [app] [schema]",
			Action: func(c *cli.Context) error {
				if os.Args[2] != "" && os.Args[3] != "" {
					err := execute(os.Args[2], "", os.Args[3], postgres.SeedPostGresWithMongo, pgtables)

					if err != nil {
						panic(err)
					}
				}

				return nil
			},
		},
		{
			// To run "go run *.go ms"
			Name:    "mongo-seed",
			Aliases: []string{"ms"},
			Usage:   "command -- ms [app] [schema]",
			Action: func(c *cli.Context) error {
				if os.Args[2] != "" && os.Args[3] != "" {
					err := execute(os.Args[2], "", os.Args[3], mongo.SeedMongoWithPostGres, mstables)

					if err != nil {
						panic(err)
					}
				}

				return nil
			},
		},
		{
			// To run "go run *.go cm"
			Name:    "create-master",
			Aliases: []string{"cm"},
			Usage:   "command -- cm [app] [schema]",
			Action: func(c *cli.Context) error {
				if os.Args[2] != "" && os.Args[3] != "" {
					var pgdb postgres.Postgres
					pgdb.Init(os.Args[2])
					err := master.CreateMaster(os.Args[2], "", os.Args[3], pgdb)

					if err != nil {
						panic(err)
					}
				}

				return nil
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func execute(app string, accountNumber string, scname string, f fn, tables map[string]interface{}) error {
	for _, table := range tables[scname].([]string) {
		err := f(app, fmt.Sprintf("%s-%s", accountNumber, scname), table)
		if err != nil {
			return err
		}
	}

	return nil
}
