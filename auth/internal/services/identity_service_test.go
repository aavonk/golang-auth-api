package services

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
	"github.com/todo-app/internal/identity"
	"github.com/todo-app/internal/repositories"
	"github.com/todo-app/testutil"
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

// Test user login should should return the user if login is successful
// and an error if not.
func TestHandleLoginFail(t *testing.T) {
	testutil.SetupUserTable(db, t)
	service := NewIdentityService(db)

	password := "supersecret"

	var users []domain.User

	for i := 0; i < 5; i++ {
		u, err := createTestUser(db, &repositories.UserDBModel{
			ID:             uuid.New(),
			FirstName:      "Aaron",
			LastName:       "testing",
			Email:          testutil.MakeRandEmail(),
			Password:       password,
			EmailConfirmed: false,
		}, t)

		if err != nil {
			t.Errorf("Error: %s", err)
			t.Error("Failed creating test users")
		}
		users = append(users, u)
	}

	// End Setup -- begin assertions
	for _, u := range users {
		request := identity.LoginRequest{
			Email:     u.Email,
			Passsword: "wrong password",
		}

		user, err := service.HandleLogin(&request)

		// We expect both user to be empty and an error returned
		// since the wrong password was given for each one
		if !user.IsEmpty() {
			t.Error("Handle login failed to return user")
		}

		if err == nil {
			t.Errorf("Received error: %s", err)
		}

	}

	testutil.TeardownUserTable(db, t)
}

// ---------------------  Helpers ---------------------------- //

func createTestUser(db *sqlx.DB, model *repositories.UserDBModel, t *testing.T) (domain.User, error) {
	newPass, err := identity.HashPassword([]byte(model.Password))

	if err != nil {
		t.Error("Failed to hash test users password")
	}

	model.Password = string(newPass)
	_, err = db.NamedExec(`INSERT INTO users (id, first_name, last_name, email, password, email_confirmed)
	 VALUES (:id, :first_name, :last_name, :email, :password, :email_confirmed)`, model)

	if err != nil {
		return domain.User{}, err
	}

	return model.ToDomain(), nil
}
