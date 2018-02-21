package main

import (
	"flag"
	"fmt"

	"github.com/avecost/ewallet"
)

const (
	// APIVersion what version of the API
	APIVersion = "1.0"
)

func main() {

	// define the parameters
	addr := flag.String("addr", ":8080", "address of our application")
	dbuser := flag.String("user", "postgres", "database user")
	dbpass := flag.String("pass", "p@ssw0rd", "database user password")
	dbname := flag.String("db", "inventiv_raffle", "database to use")
	dbaddr := flag.String("dbaddr", "localhost", "database address & port")
	// parse the flag
	flag.Parse()

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", *dbuser, *dbpass, *dbaddr, *dbname)
	// create a new server
	srvr := ewallet.NewServer(connStr)
	// run the server
	srvr.Run(*addr)
}
