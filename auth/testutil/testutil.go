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
	CREATE TABLE IF NOT EXISTS "users" (
		"id" TEXT NOT NULL,
		"first_name" VARCHAR(50) NOT NULL,
		"last_name" VARCHAR(80) NOT NULL,
		"email" TEXT NOT NULL UNIQUE,
		"email_confirmed" BOOLEAN NOT NULL DEFAULT false,
		"password" TEXT NOT NULL,

		PRIMARY KEY ("id")
	);
	`
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
