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
	createdUser, err := createTestUser(db, &repositories.UserDBModel{
		ID:        uuid.New(),
		FirstName: "Hello",
		LastName:  "Goodbye",
		Email:     "email@email.com",
		Password:  password,
		Activated: false,
	}, t)

	if err != nil {
		t.Fatal("failed to create test user")
	}

	tests := []struct {
		name       string
		request    *identity.LoginRequest
		shouldPass bool
	}{
		{
			name: "Incorrect email",
			request: &identity.LoginRequest{
				Email:     "randomemail@email.com",
				Passsword: password,
			},
			shouldPass: false,
		},
		{
			name: "Correct email, wrong password",
			request: &identity.LoginRequest{
				Email:     createdUser.Email,
				Passsword: password,
			},
			shouldPass: false,
		},
		{
			name: "Correct Email & Password",
			request: &identity.LoginRequest{
				Email:     "email@email.com",
				Passsword: password,
			},
			shouldPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := service.HandleLogin(tt.request)

			if err != nil && tt.shouldPass {
				t.Errorf("Failedl login with error %v", err)
			}

			if u == nil && tt.shouldPass {
				t.Errorf("Expected a user returned, but got isEmpty")
			}
		})
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
		Name                 string
	}
	items := []options{
		{
			ShouldFail:           true,
			CreateUserBeforeTest: false,
			Name:                 "Bad password",
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
			Name:                 "Invalid Email",
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
			Name:                 "Valid",
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
			Name:                 "Existing user",
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
		t.Run(item.Name, func(t *testing.T) {
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
			if (item.ShouldFail && err == nil) || (item.ShouldFail && user != nil) {
				t.Errorf("Item was supposed to fail but did not. Reason that it was supposed to fail: %s", item.Reason)
			}

			// If the test is supposed to pass and there is an error, or the user returned is empty
			// then something went wrong.
			if (!item.ShouldFail && err != nil) || (!item.ShouldFail && user == nil) {
				t.Errorf("Was not supposed to fail but did. Error: %s. Reason supposed to fail %s", err, item.Reason)
			}

		})

	}

	testutil.TeardownUserTable(db, t)

}

func TestGetUserById(t *testing.T) {
	testutil.SetupUserTable(db)
	service := NewIdentityService(db)

	idOne := uuid.New()
	idTwo := uuid.New()

	tests := []struct {
		name     string
		userId   string
		wantUser domain.User
		wantErr  error
	}{
		{
			name:   "Should pass 1",
			userId: idOne.String(),
			wantUser: domain.User{
				ID:        idOne,
				FirstName: "Hello",
				LastName:  "Goodbye",
				Email:     "test1@gmail.com",
				Password:  "hellohello",
				Activated: false,
			},
			wantErr: nil,
		},
		{
			name:   "should pass 2",
			userId: idTwo.String(),
			wantUser: domain.User{
				ID:        idTwo,
				FirstName: "Goodbye",
				LastName:  "Hello",
				Email:     "email3@gmail.com",
				Password:  "asdfalsdkf",
				Activated: false,
			},
			wantErr: nil,
		},
		{
			name:   "Should throw error",
			userId: "fakeIdthat0123kjf",
			wantUser: domain.User{
				ID:        uuid.New(),
				FirstName: "hellloooo",
				LastName:  "goooodbye",
				Email:     "asdf@gmail.com",
				Password:  "asdfasdf",
				Activated: false,
			},
			wantErr: repositories.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createTestUser(db, &repositories.UserDBModel{
				ID:        tt.wantUser.ID,
				FirstName: tt.wantUser.FirstName,
				LastName:  tt.wantUser.LastName,
				Email:     tt.wantUser.Email,
				Password:  tt.wantUser.Password,
				Activated: tt.wantUser.Activated,
			}, t)

			_, err := service.GetUserById(tt.userId)

			if err != tt.wantErr {
				t.Errorf("want: %v; got %v", tt.wantErr, err)
			}

		})
	}

	testutil.TeardownUserTable(db, t)
}

// ---------------------  Helpers ---------------------------- //

func createTestUser(db *sqlx.DB, model *repositories.UserDBModel, t *testing.T) (*domain.User, error) {
	newPass, err := identity.HashPassword([]byte(model.Password))

	if err != nil {
		t.Error("Failed to hash test users password")
	}

	model.Password = string(newPass)
	_, err = db.NamedExec(`INSERT INTO users (id, first_name, last_name, email, password, activated)
	 VALUES (:id, :first_name, :last_name, :email, :password, :activated)`, model)

	if err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}
