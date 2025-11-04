package postgres

import (
	"database/sql"

	// postgres driver
	_ "github.com/lib/pq"
	"github.com/lopezator/migrator"
)

// defaultOpts is the migration options
var defaultOpts = []migrator.Option{migrator.WithLogger(newNopLogger())}

type nopLogger struct{}

func (*nopLogger) Printf(string, ...interface{}) {}

func newNopLogger() migrator.Logger { return &nopLogger{} }

// Migrate migrates the database.
func Migrate(db *sql.DB, opts ...migrator.Option) error {
	if len(opts) == 0 {
		opts = defaultOpts
	}

	opts = append(opts, migrations)

	m, err := migrator.New(opts...)
	if err != nil {
		return err
	}

	return m.Migrate(db)
}

// MustMigrate migrates the database and panics if an error occurs.
func MustMigrate(db *sql.DB, opts ...migrator.Option) {
	if len(opts) == 0 {
		opts = defaultOpts
	}

	opts = append(opts, migrations)

	m, err := migrator.New(opts...)
	if err != nil {
		panic(err)
	}

	err = m.Migrate(db)
	if err != nil {
		panic(err)
	}
}

var migrations = migrator.Migrations(
	&migrator.Migration{
		Name: "Create cards table",
		Func: func(tx *sql.Tx) error {
			stmnt := `CREATE TABLE IF NOT EXISTS "cards" (
				id varchar(32) PRIMARY KEY,
				external_id varchar(36),
				deleted_at timestamp,
				created_at timestamp NOT NULL DEFAULT now()
			)`
			if _, err := tx.Exec(stmnt); err != nil {
				return err
			}

			return nil
		},
	},
)
