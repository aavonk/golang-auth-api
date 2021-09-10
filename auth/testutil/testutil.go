package testutil

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

func SetupUserTable(db *sqlx.DB) {
	var schema = `
	CREATE EXTENSION IF NOT EXISTS citext;

	CREATE TABLE IF NOT EXISTS users (
		id text NOT NULL PRIMARY KEY,
		created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
		first_name text NOT NULL,
		last_name text NOT NULL,
		email citext UNIQUE NOT NULL,
		password bytea NOT NULL,
		activated bool NOT NULL
	);`
	// log.Println("**** Creating User Table ****")
	db.MustExec(schema)

}

// Removes the table from the test db
func TeardownUserTable(db *sqlx.DB, t *testing.T) {
	_, err := db.Exec(`DROP TABLE IF EXISTS "users"`)
	if err != nil {
		t.Error("Failed to clear user table")
	}
	// log.Println("**** Successfully Dropped User Table ****")

}

func MakeRandEmail() string {
	b := make([]byte, 10)
	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return fmt.Sprintf("%s@gmail.com", string(b))
}
