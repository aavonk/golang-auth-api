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

// TestHandleLogin should should return the user if login is successful
// and an error if not.
func TestHandleLogin(t *testing.T) {
	testutil.SetupUserTable(db)
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
	// Check that the service returns an empty user and an error
	// given the wrong password
	for _, u := range users {
		request := identity.LoginRequest{
			Email:     u.Email,
			Passsword: "wrong password",
		}

		user, err := service.HandleLogin(&request)

		// We expect both user to be empty and an error returned
		// since the wrong password was given for each one
		if !user.IsEmpty() {
			t.Error("HandleLogin incorrectly returned a user given the wrong password")
		}

		if err == nil {
			t.Errorf("Received error while logging in with incorrect password: %v", err)
		}

	}

	// Check that the service returns the user object and no error
	// given the correct email/password
	for _, u := range users {
		request := identity.LoginRequest{
			Email:     u.Email,
			Passsword: password,
		}

		user, err := service.HandleLogin(&request)

		if user.IsEmpty() {
			t.Errorf("HandleLogin returned empty user. Expected: %+v", u)
		}

		if err != nil {
			t.Errorf("Received error while logging in with correct password: %v", err)
		}
	}

	testutil.TeardownUserTable(db, t)
}

// TestHandleRegistration should return the user object is registration
// is successful and an error if not. There are many points where registration might fail e.g.
// - request does not pass model validation (password requirement, invalid email, etc.)
// - There is already an existing account with the requested email
// - There is an error inserting the record into the database
func TestHandleRegistration(t *testing.T) {
	testutil.SetupUserTable(db)
	service := NewIdentityService(db)
	// password := "password"

	type options struct {
		ShouldFail           bool
		Reason               string
		User                 domain.User
		CreateUserBeforeTest bool
	}
	items := []options{
		{
			ShouldFail:           true,
			CreateUserBeforeTest: false,
			Reason:               "[Model Validation]-Password doesn't meet requirements",
			User: domain.User{
				ID:        uuid.New(),
				FirstName: "Test",
				LastName:  "Testing",
				Email:     testutil.MakeRandEmail(),
				Password:  "1234",
			},
		},
		{
			ShouldFail:           true,
			CreateUserBeforeTest: false,
			Reason:               "[Model Validation]-Invalid email provided",
			User: domain.User{
				ID:        uuid.New(),
				FirstName: "Test",
				LastName:  "Test",
				Email:     "email.com",
				Password:  "password",
			},
		},
		{
			ShouldFail:           false,
			CreateUserBeforeTest: false,
			Reason:               "Should not fail",
			User: domain.User{
				ID:        uuid.New(),
				FirstName: "Test",
				LastName:  "Test",
				Email:     "test@testing.com",
				Password:  "longpassword",
			},
		},
		{
			ShouldFail:           true,
			CreateUserBeforeTest: true,
			Reason:               "[Existing user] - User with same email already exists. Should fail",
			User: domain.User{
				ID:        uuid.New(),
				FirstName: "Test",
				LastName:  "Test",
				Email:     "aaron@testing.com",
				Password:  "password",
			},
		},
	}

	for _, item := range items {
		// If we create a user before the actua test implementation is passed,
		// the test case should fail because there will be an existing user in the
		// database with the same email. This creates that user, and then allows
		// the rest of the test case to continue.
		if item.CreateUserBeforeTest {
			_, err := service.HandleRegister(&item.User)
			if err != nil {
				t.Errorf("Failed creating user before tests. Error: %v", err)
			}
		}

		user, err := service.HandleRegister(&item.User)

		// If the test case is supposed to fail, and no error is present
		// something went wrong
		if (item.ShouldFail && err == nil) || (item.ShouldFail && !user.IsEmpty()) {
			t.Errorf("Item was supposed to fail but did not. Reason that it was supposed to fail: %s", item.Reason)
		}

		// If the test is supposed to pass and there is an error, or the user returned is empty
		// then something went wrong.
		if (!item.ShouldFail && err != nil) || (!item.ShouldFail && user.IsEmpty()) {
			t.Errorf("Was not supposed to fail but did. Error: %s. Reason supposed to fail %s", err, item.Reason)
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
