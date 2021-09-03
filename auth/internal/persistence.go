package internal

import (
	_ "database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DataStore struct {
	Client *sqlx.DB
}

// GetDataStore connects to a postgrest database, makes sure it has a valid
// conenction by calling Ping(), and then returns a pointer to the db
// client if successful, otherwise it returns an error.
func GetDataStore(connStr string) (*DataStore, error) {
	db, err := getDatabaseConn(connStr)

	if err != nil {
		return nil, err
	}

	return &DataStore{
		Client: db,
	}, nil
}

func (d *DataStore) Close() error {
	return d.Client.Close()
}

func getDatabaseConn(connString string) (*sqlx.DB, error) {
	// Connect calls sql.Ping() internally
	db, err := sqlx.Connect("postgres", connString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
