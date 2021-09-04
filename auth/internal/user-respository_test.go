package internal

import (
	_ "database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/todo-app/internal/domain"
)

var db *sqlx.DB

// TestMain allows us to spin up a new docker container and connection
// to a postgres database for each test.
// It can only be called once in the same package.
func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it.
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_PASSWORD=password",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=testdb",
		},
	}, func(config *docker.HostConfig) {
		// Set autoremove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://postgres:password@%s/testdb?sslmode=disable", hostAndPort)
	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accep connections yet.
	pool.MaxWait = 120 * time.Second
	if err := pool.Retry(func() error {
		db, err = sqlx.Connect("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)

	}
	os.Exit(code)
}
func TestUserRepositoryGetByEmail(t *testing.T) {
	createUserTable(db, t)

	users := []UserDBModel{
		{
			ID:             uuid.New(),
			FirstName:      "Test",
			LastName:       "NumberOne",
			Email:          "test1@gmail.com",
			Password:       "testing",
			EmailConfirmed: false,
		},
		{
			ID:             uuid.New(),
			FirstName:      "Test",
			LastName:       "NumberTwo",
			Email:          "test2@gmail.com",
			Password:       "testing2",
			EmailConfirmed: false,
		},
		{
			ID:             uuid.New(),
			FirstName:      "Test",
			LastName:       "NumberThree",
			Email:          "test3@gmail.com",
			Password:       "testing3",
			EmailConfirmed: false,
		},
	}

	for _, user := range users {
		created, err := createTestUser(db, &user)

		if err != nil {
			t.Error("Failed to create test user")
		}
		repository := NewUserRepository(db)

		// See if the user we previously created exists
		existing := repository.GetByEmail(created.Email)

		// if existing is an empty user struct, then we didn't
		// find the user we just created, and it should fail
		if existing.IsEmpty() == true {
			t.Error("created user not found")
			t.Errorf("existing: %+v \n", existing)
			t.Errorf("created: %+v \n", created)
		}

		if created != existing {
			t.Errorf("Created is not the same as existing. created %+v, existing: %+v", created, existing)
		}
	}

	clearUserTable(db, t)
}

// Helpers

// createTestUser inserts into the database directly and does not
// hash a password. Do NOT use this when testing the Create method on
// the user repository
func createTestUser(db *sqlx.DB, model *UserDBModel) (domain.User, error) {

	_, err := db.NamedExec(`INSERT INTO users (id, first_name, last_name, email, password, email_confirmed)
	 VALUES (:id, :first_name, :last_name, :email, :password, :email_confirmed)`, model)

	if err != nil {
		return domain.User{}, err
	}

	return model.ToDomain(), nil
}

// Applies the user schema to the test db
func createUserTable(db *sqlx.DB, t *testing.T) {
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
	log.Println("**** Creating User Table ****")
	db.MustExec(schema)
}

// Removes the table from the test db
func clearUserTable(db *sqlx.DB, t *testing.T) {
	_, err := db.Exec(`DROP TABLE IF EXISTS "users"`)
	if err != nil {
		t.Error("Failed to clear user table")
	}
	log.Println("**** Successfully Dropped User Table ****")

}
