package repository_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	repo "user-service/interface/repository"
	"user-service/models"
	"user-service/usecases/repository"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	port     = "5435"
	user     = "postgres"
	password = "postgres"
	dbName   = "users_test"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=20"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDb *sql.DB
var userRepo repository.UserRepository

func TestMain(m *testing.M) {
	// connect to docker; fail if docker not running
	p, err := dockertest.NewPool("")

	if err != nil {
		log.Fatalf("could not connect to docker, is it running? %s", err)
	}

	pool = p

	// setup docker options, specifying the image and so forth
	opt := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15.1-alpine",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// get a resource (docker image)
	resource, err = pool.RunWithOptions(&opt)
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	// start the image and wait until its ready
	if err := pool.Retry(func() error {
		var err error
		testDb, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("error:", err)
		}
		return testDb.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database: %s", err)
	}

	// populate database with empty table
	err = createTables()
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("cant create table: %s", err)
	}

	userRepo = repo.NewUserRepository(testDb)
	code := m.Run()

	// clean up
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("cant clean resources: %s", err)
	}

	os.Exit(code)
}

func Test_userRepo_pingDB(t *testing.T) {
	if err := testDb.Ping(); err != nil {
		t.Errorf("cant ping database")
	}
}

func Test_userRepo_Create(t *testing.T) {
	testUser := models.UserPayload{
		Fname:    "ryan",
		Lname:    "pujo",
		Username: "ryanpujo",
		Email:    "ryanpujo",
		Password: "secret",
	}

	id, err := userRepo.Create(&testUser)
	if err != nil {
		t.Errorf("create user returned an error: %s", err)
	}

	if id != 1 {
		t.Errorf("create user returned a wrong id; expected 1 but got %d", id)
	}
}

func Test_userRepo_FindByid(t *testing.T) {
	result, err := userRepo.FindById(int64(1))
	if err != nil {
		t.Errorf("find by id returned an error: %s", err)
	}
	if int64(result.Id) != 1 {
		t.Errorf("create user returned a wrong id; expected 1 but got %d", result.Id)
	}
	if result.Fname != "ryan" {
		t.Errorf("create user returned a wrong first name; expected ryan but got %s", result.Fname)
	}
}

func Test_userRepo_FindByUsername(t *testing.T) {
	result, err := userRepo.FindByUsername("ryanpujo")
	if err != nil {
		t.Errorf("find by id returned an error: %s", err)
	}
	if result.Username != "ryanpujo" {
		t.Errorf("create user returned a wrong id; expected ryanpujo but got %s", result.Username)
	}
}

func Test_userRepo_FindUsers(t *testing.T) {
	users, err := userRepo.FindUsers()
	if err != nil {
		t.Errorf("find users returned an error: %s", err)
	}
	if len(users) != 1 {
		t.Errorf("users reported a wrong size; expected 1, but got %d", len(users))
	}
}

func Test_userRepo_UpdateById(t *testing.T) {
	testUser := models.UserPayload{
		Id:       1,
		Fname:    "ryan",
		Lname:    "conor",
		Username: "ryanpujo",
		Email:    "ryanpujo",
		Password: "secret",
	}
	err := userRepo.Update(&testUser)
	if err != nil {
		t.Errorf("failed to update user with id 1: %s", err)
	}

	user, _ := userRepo.FindById(1)
	if user.Lname != "conor" {
		t.Errorf("last name update failed expect conor but got %s", user.Lname)
	}
}

func Test_userRepo_DeleteById(t *testing.T) {
	err := userRepo.DeleteById(1)
	if err != nil {
		t.Errorf("failed to deleste user with id 1: %s", err)
	}

	_, err = userRepo.FindById(1)
	if err == nil {
		t.Errorf("retrieved a user that should have been deleted")
	}
}

func createTables() (err error) {
	var tableSql []byte
	tableSql, err = os.ReadFile("./testdata/users.sql")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = testDb.Exec(string(tableSql))
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
