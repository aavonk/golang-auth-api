package internal

import (
	"context"
	_ "database/sql"
	"time"

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
	db, err := sqlx.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) conenctions in the poo;.
	// Passing a value less than or equal to 0 will mean ther is no limit
	db.SetMaxOpenConns(25)

	// Set the maximum number of idle connections in the pool. Again,
	// anything <= 0 will mean ther is no limit.
	// Note: It should always be less than or equal to the Max Open Connections
	db.SetMaxIdleConns(25)

	duration, err := time.ParseDuration("15m")
	if err != nil {
		return nil, err
	}
	// Set the maximum idle timeout
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
