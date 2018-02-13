package main

import "fmt"

const (
	// APIVersion what version of the API
	APIVersion = "1.0"

	// DBUser user account of the DB
	DBUser = "postgres"
	// DBPass password of the DB user
	DBPass = "p@ssw0rd"
	// DBName name of the Database
	DBName = "inventiv_raffle"
)

func main() {

	conn, err := Connect(DBUser, DBPass, DBName)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("PostgreSQL Connect Succesful")

	h := NewHandler(conn)

	h.Run()
}
