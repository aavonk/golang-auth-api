package config

import (
	"flag"
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
}

func Get() *Confg {
	c := &Confg{}

	flag.StringVar(&c.dbUser, "dbuser", os.Getenv("POSTGRES_USER"), "Database username")
	flag.StringVar(&c.dbPassword, "dbpassword", os.Getenv("POSTGRES_PASSWORD"), "Database password")
	flag.StringVar(&c.dbHost, "dbhost", os.Getenv("POSTGRES_HOST"), "Database host")
	flag.StringVar(&c.dbPort, "dbport", os.Getenv("POSTGRES_PORT"), "Database port")
	flag.StringVar(&c.dbName, "dbname", os.Getenv("DB_NAME"), "Database name")
	flag.StringVar(&c.testDBHost, "testdbhost", os.Getenv("TEST_DB_HOST"), "Test database host")
	flag.StringVar(&c.testDBName, "testdbname", os.Getenv("TEST_DB_NAME"), "Test db name")
	flag.StringVar(&c.apiPort, "apiport", os.Getenv("API_PORT"), "Api port to listen on")

	flag.Parse()

	return c
}
