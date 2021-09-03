package internal

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DataStore struct {
	Client *sql.DB
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

func getDatabaseConn(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
