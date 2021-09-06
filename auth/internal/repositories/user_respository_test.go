package repositories

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

// TestGetById tests whether the GetById method of the user-repository successfully returns
// a user domain object when a user is found, and an empty domain object when no record is found
func TestGetById(t *testing.T) {
	testutil.SetupUserTable(db)
	repository := NewUserRepository(db)

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

	// Generate users in the DB
	for _, u := range users {
		_, err := CreateTestUser(db, u)
		if err != nil {
			t.Errorf("Failed to create test users. Error: %v", err)
		}
	}

	// Implement test
	for _, user := range users {
		found := repository.GetById(user.ID.String())

		if found.IsEmpty() {
			t.Errorf("Failed to find user. Expected: %+v \t Got: %+v", user, found)
		}
	}

	testutil.TeardownUserTable(db, t)
}

// TestGetByemail tests whether the GetByEmail method of the user-repository correctly
// returns a user when found from the database with a certain email
func TestGetByEmail(t *testing.T) {
	testutil.SetupUserTable(db)
	repository := NewUserRepository(db)

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
		created, err := CreateTestUser(db, user)

		if err != nil {
			t.Error("Failed to create test user")
		}

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

	testutil.TeardownUserTable(db, t)
}

// TestGetByEmailNotFound tests whether the GetByEmail method of the user-repository
// successfully returns an empty domain User object when no user is found by the given email
func TestGetByEmailNotFound(t *testing.T) {
	// Apply the user schema to the db
	testutil.SetupUserTable(db)
	repository := NewUserRepository(db)

	emails := []string{"aaron@gmail.com", "testing@test.com", "test1@gmail.com", "hello@test.com"}

	for _, email := range emails {

		// Since the DB will be empty, we should expect no users to be found no matter what
		// email we provide
		found := repository.GetByEmail(email)

		if found.IsEmpty() != true {
			t.Errorf("Found a user and returned. Got %+v, want %+v", found, domain.User{})
		}
	}
	// Drop the user table for good measure
	testutil.TeardownUserTable(db, t)
}

// TestUserCreate tests whether a user is successfully saved to the database
// given a user domain object, and returns a user domain object.
func TestUserCreate(t *testing.T) {

	testutil.SetupUserTable(db)
	repository := NewUserRepository(db)

	users := []domain.User{
		{
			ID:        uuid.New(),
			FirstName: "Aaron",
			LastName:  "von Kreisler",
			Email:     "aaron@email.com",
			Password:  "password",
		},
		{
			ID:        uuid.New(),
			FirstName: "Billy",
			LastName:  "Bob",
			Email:     "billybob@email.com",
			Password:  "password",
		},
		{
			ID:        uuid.New(),
			FirstName: "Test",
			LastName:  "Testing",
			Email:     "hello@testing.com",
			Password:  "password",
		},
	}

	for i, user := range users {
		u, err := repository.Create(&user)

		if err != nil {
			t.Errorf("Failed to create user: %s", err)
		}

		if u.IsEmpty() {
			t.Error("Returned an empty user object")
		}

		if u.ID != users[i].ID {
			t.Errorf("Given userID does not match with received. Got %s want %s", u.ID, users[i].ID)
		}

		if u.FirstName != users[i].FirstName {
			t.Errorf("Given FirstName does not match received FirstName. Got %s want %s", u.FirstName, users[i].FirstName)
		}

		if u.LastName != users[i].LastName {
			t.Errorf("Given LastName does not match received LastName. Got %s want %s", u.LastName, users[i].LastName)

		}
		if u.Email != users[i].Email {
			t.Errorf("Given Email does not match received Email. Got %s want %s", u.Email, users[i].Email)

		}

		// This check is different because if we give a password, it should NOT return the actual password,
		// but rather it should return the encrypted version. Ex: given password of "password", the user returned from create
		// should have a password of something like "asdf0978wklj2340usdfjhsf08734kjhsdf8".
		if u.Password == users[i].Password {
			t.Error("FAILED TO ENCRYPT USERS PASSWORD")
		}
	}

	testutil.TeardownUserTable(db, t)
}

// TestUserCreateFailed tests whether the Create method handles invalid data
// correctly and returns an empty user domain object along with the error
func TestUserCreateFailed(t *testing.T) {
	// We can simulate an error by NOT setting up the User table.
	// By doing so, the db query will fail to execute and should return an error
	repository := NewUserRepository(db)
	users := []domain.User{
		{
			ID:        uuid.New(),
			FirstName: "Aaron",
			LastName:  "von Kreisler",
			Email:     "aaron@email.com",
			Password:  "password",
		},
		{
			ID:        uuid.New(),
			FirstName: "Billy",
			LastName:  "Bob",
			Email:     "billybob@email.com",
			Password:  "password",
		},
		{
			ID:        uuid.New(),
			FirstName: "Test",
			LastName:  "Testing",
			Email:     "hello@testing.com",
			Password:  "password",
		},
	}

	for _, user := range users {
		u, err := repository.Create(&user)

		if err == nil {
			t.Error("Expected an error, did not receive one")
		}

		if !u.IsEmpty() {
			t.Errorf("Expected an empty user, received %+v", u)
		}
	}
}

// ---------------------  Helpers ---------------------------- //
// CreateTestUser is a helper function that inserts a user to the DB given a UserDBModel.
// Note: It does NOT hash passwords.
func CreateTestUser(db *sqlx.DB, model UserDBModel) (domain.User, error) {

	_, err := db.NamedExec(`INSERT INTO users (id, first_name, last_name, email, password, email_confirmed)
	 VALUES (:id, :first_name, :last_name, :email, :password, :email_confirmed)`, model)

	if err != nil {
		return domain.User{}, err
	}

	return model.ToDomain(), nil
}
