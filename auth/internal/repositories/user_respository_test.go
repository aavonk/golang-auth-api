package repositories

import (
	_ "database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
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
	repo := NewUserRepository(db)

	tests := []struct {
		User       UserDBModel
		Name       string
		ShouldFail bool
	}{
		{
			Name:       "One.Should Pass",
			ShouldFail: false,
			User: UserDBModel{
				ID:        uuid.New(),
				FirstName: "test",
				LastName:  "test",
				Email:     "test1@gmail.com",
				Password:  "password",
				Activated: false,
				CreatedAt: time.Time{},
			},
		},
		{
			Name:       "Two.Should Pass",
			ShouldFail: false,
			User: UserDBModel{
				ID:        uuid.New(),
				FirstName: "test",
				LastName:  "test",
				Email:     "test2@gmail.com",
				Password:  "password",
				Activated: false,
				CreatedAt: time.Time{},
			},
		},
		{
			Name:       "Three.Should Fail",
			ShouldFail: true,
			User: UserDBModel{
				ID:        uuid.New(),
				FirstName: "test",
				LastName:  "test",
				Email:     "test3@gmail.com",
				Password:  "password",
				Activated: false,
				CreatedAt: time.Time{},
			},
		},
		{
			Name:       "Four.Should Fail",
			ShouldFail: true,
			User: UserDBModel{
				ID:        uuid.New(),
				FirstName: "test",
				LastName:  "test",
				Email:     "test4@gmail.com",
				Password:  "password",
				Activated: false,
				CreatedAt: time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			var userId string
			// If we want the test to fail (marked by ShouldFail) then we
			// can pass in a UserID that we know will not exist in the DB,
			// as all UserIDS are UUID's in the form of 000-000-00-00-0-...
			// Otherwise, just keep the UserID generated by UUID
			//This will cause the repo.GetById method to fail and we can test what is expected

			if !tt.ShouldFail {
				userId = tt.User.ID.String()
			} else {
				userId = "asdfas"
			}

			// Create the user using the actual uuid
			_, err := CreateTestUser(db, tt.User)
			if err != nil {
				t.Errorf("error creating user: %v", err)
			}

			// pass the dynamically made userId to the wuery
			_, err = repo.GetById(userId)

			if err != nil && !tt.ShouldFail {
				t.Errorf("Was not supposed to fail but did. err %v", err)
			}

			// If we are supposed to fail and there is an error,
			// check to see if it is an ErrRecordNotFound error
			if tt.ShouldFail && err != nil {
				if !errors.Is(err, ErrRecordNotFound) {
					t.Errorf("Received unexpected err. Expected ErrRecordNotFound, got: %v", err)
				}
			}

		})
	}
	testutil.TeardownUserTable(db, t)

}

// TestGetByEmail will test to see if it successfully returns
// The user when a user is found and the correct user when one isn't found.
func TestGetByEmail(t *testing.T) {
	testutil.SetupUserTable(db)
	repo := NewUserRepository(db)

	tests := []struct {
		User       UserDBModel
		Name       string
		ShouldFail bool
	}{
		{
			Name:       "One.Should Pass",
			ShouldFail: false,
			User: UserDBModel{
				ID:        uuid.New(),
				FirstName: "test",
				LastName:  "test",
				Email:     "test1@gmail.com",
				Password:  "password",
				Activated: false,
				CreatedAt: time.Time{},
			},
		},
		{
			Name:       "Two.Should Pass",
			ShouldFail: false,
			User: UserDBModel{
				ID:        uuid.New(),
				FirstName: "test",
				LastName:  "test",
				Email:     "test2@gmail.com",
				Password:  "password",
				Activated: false,
				CreatedAt: time.Time{},
			},
		},
		{
			Name:       "Three.Should Fail",
			ShouldFail: true,
			User: UserDBModel{
				ID:        uuid.New(),
				FirstName: "test",
				LastName:  "test",
				Email:     "test3@gmail.com",
				Password:  "password",
				Activated: false,
				CreatedAt: time.Time{},
			},
		},
		{
			Name:       "Four.Should Fail",
			ShouldFail: true,
			User: UserDBModel{
				ID:        uuid.New(),
				FirstName: "test",
				LastName:  "test",
				Email:     "test4@gmail.com",
				Password:  "password",
				Activated: false,
				CreatedAt: time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			var email string

			// If we want the test to fail (marked by ShouldFail) then we
			// can pass in an email that we know will not exist in the DB,
			// as the only users that will exist in the database are the ones
			// listed in the test items. This will cause the repo.GetByEmail
			// method to fail and we can test what is expected
			if tt.ShouldFail {
				email = "fake@email.com"
			} else {
				email = tt.User.Email
			}
			// create a user with a legit email
			_, err := CreateTestUser(db, tt.User)
			if err != nil {
				t.Errorf("Failed creating user before test: %v", err)
			}

			// Run the test with the generated emaill var
			_, err = repo.GetByEmail(email)

			if tt.ShouldFail && err == nil {
				t.Errorf("Expected to fail but did not")
			}

			if tt.ShouldFail && err != nil {
				if !errors.Is(err, ErrRecordNotFound) {
					t.Errorf("Expected ErrRecordNotFound got: %v", err)
				}
			}

			if !tt.ShouldFail && err != nil {
				t.Errorf("Was not suppossed to fail but did. Err: %v", err)
			}

		})
	}

	testutil.TeardownUserTable(db, t)
}

func TestCreate(t *testing.T) {
	testutil.SetupUserTable(db)
	repo := NewUserRepository(db)

	tests := []struct {
		Name    string
		WantErr error
		User    *domain.User
	}{
		{
			Name:    "Valid 1",
			WantErr: nil,
			User: &domain.User{
				ID:        uuid.New(),
				FirstName: "test",
				LastName:  "test",
				Email:     "test@gmail.com",
				Password:  "asdfasdf",
				Activated: false,
			},
		},
		{
			Name:    "Duplicate Email",
			WantErr: ErrDuplicateEmail,
			User: &domain.User{
				ID:        uuid.New(),
				FirstName: "test",
				LastName:  "test",
				Email:     "test@gmail.com",
				Password:  "hello",
				Activated: false,
			},
		},
		{
			Name:    "Valid 2",
			WantErr: nil,
			User: &domain.User{
				ID:        uuid.New(),
				FirstName: "hello",
				LastName:  "goodbye",
				Email:     "email777@gmail.com",
				Password:  "passwordthatslong",
				Activated: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := repo.Create(tt.User)

			if err != tt.WantErr {
				t.Errorf("want: %v, got %v", tt.WantErr, err)
				t.FailNow()
			}
		})
	}

	testutil.TeardownUserTable(db, t)
}

func TestUpdate(t *testing.T) {
	testutil.SetupUserTable(db)
	repo := NewUserRepository(db)

	tests := []struct {
		Name             string
		WantErr          error
		User             *domain.User
		CreateBeforeTest bool
	}{
		{
			Name:    "Success1",
			WantErr: nil,
			User: &domain.User{
				ID:        uuid.New(),
				FirstName: "test",
				LastName:  "test",
				Email:     "test@gmail.com",
				Password:  "asdfasdf",
				Activated: false,
			},
		},
		{
			Name:    "Success 2",
			WantErr: nil,
			User: &domain.User{
				ID:        uuid.New(),
				FirstName: "test",
				LastName:  "test",
				Email:     "test7676@gmail.com",
				Password:  "hello",
				Activated: false,
			},
		},
		{
			Name:    "Success 3",
			WantErr: nil,
			User: &domain.User{
				ID:        uuid.New(),
				FirstName: "hello",
				LastName:  "goodbye",
				Email:     "email777@gmail.com",
				Password:  "passwordthatslong",
				Activated: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			_, err := repo.Create(tt.User)
			if err != nil {
				t.Errorf("failed creating user before test: %s", err)
			}

			// Make a change to the user *these are pointers, so the actual memory location value will be changed*
			tt.User.Activated = true
			tt.User.LastName = "Test Lastname"

			err = repo.Update(tt.User)

			if err != tt.WantErr {
				t.Errorf("wanted %s; got %s", tt.WantErr, err)
			}

			// Create a new user struct with the updates we expect to be reflected
			want := &domain.User{
				ID:        tt.User.ID,
				FirstName: tt.User.FirstName,
				LastName:  "Test Lastname",
				Email:     tt.User.Email,
				Password:  tt.User.Password,
				Activated: true,
				CreatedAt: tt.User.CreatedAt,
			}

			if !reflect.DeepEqual(tt.User, want) {
				t.Errorf("got %+v, want: %+v", tt.User, want)
			}

		})
	}
	testutil.TeardownUserTable(db, t)
}

// ---------------------  Helpers ---------------------------- //
// CreateTestUser is a helper function that inserts a user to the DB given a UserDBModel.
// Note: It does NOT hash passwords.
func CreateTestUser(db *sqlx.DB, model UserDBModel) (*domain.User, error) {

	_, err := db.NamedExec(`INSERT INTO users (id, first_name, last_name, email, password, activated)
	 VALUES (:id, :first_name, :last_name, :email, :password, :activated)`, model)

	if err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}
