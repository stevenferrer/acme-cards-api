package main

import (
	"database/sql"
	"log"
	"os"
	"sort"
	"time"

	_ "github.com/lib/pq"
	"github.com/lopezator/migrator"
	"github.com/urfave/cli/v2"

	"github.com/stevenferrer/acme-cards-api/acme/postgres"
)

func main() {
	app := &cli.App{
		Name:  "migrate",
		Usage: "db migration tool",
		Commands: []*cli.Command{
			{
				Name:  "up",
				Usage: "migrate up",
				Action: func(c *cli.Context) error {
					dsn := os.Getenv("DSN")

					db, err := sql.Open("postgres", dsn)
					if err != nil {
						return err
					}
					defer db.Close()

					err = db.Ping()
					if err != nil {
						return err
					}

					l := log.New(os.Stdout, "", 0)
					l.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + " [migrate] ")
					if err := postgres.Migrate(db, migrator.WithLogger(l)); err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
