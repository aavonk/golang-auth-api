package config

import (
	"flag"
	"fmt"
	"os"
)

type Confg struct {
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     string
	dbName     string
	testDBHost string
	testDBName string
	apiPort    string
	migrate    string
}

func Get() *Confg {
	c := &Confg{}

	flag.StringVar(&c.dbUser, "dbuser", os.Getenv("POSTGRES_USER"), "Database username")
	flag.StringVar(&c.dbPassword, "dbpassword", os.Getenv("POSTGRES_PASSWORD"), "Database password")
	flag.StringVar(&c.dbHost, "dbhost", os.Getenv("POSTGRES_HOST"), "Database host")
	flag.StringVar(&c.dbPort, "dbport", os.Getenv("POSTGRES_PORT"), "Database port")
	flag.StringVar(&c.dbName, "dbname", os.Getenv("POSTGRES_DB"), "Database name")
	flag.StringVar(&c.testDBHost, "testdbhost", os.Getenv("TEST_DB_HOST"), "Test database host")
	flag.StringVar(&c.testDBName, "testdbname", os.Getenv("TEST_DB_NAME"), "Test db name")
	flag.StringVar(&c.apiPort, "apiport", os.Getenv("API_PORT"), "Api port to listen on")
	flag.StringVar(&c.migrate, "migrate", "up", "Direction to migrate DB [up or down]")
	flag.Parse()

	return c
}

func (c *Confg) GetMigration() string {
	return c.migrate
}

func (c *Confg) getDBConnStr(dbhost, dbname string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.dbUser,
		c.dbPassword,
		dbhost,
		c.dbPort,
		dbname,
	)
}

func (c *Confg) GetDBConnStr() string {
	return c.getDBConnStr(c.dbHost, c.dbName)
}

func (c *Confg) GetTestDBConnStr() string {
	return c.getDBConnStr(c.testDBHost, c.testDBName)
}

func (c *Confg) GetAPIPort() string {
	return ":" + c.apiPort
}
