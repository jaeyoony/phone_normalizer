package sequel

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "jy"
	password = "poop"
	dbname = "golang_phone_db"
)

func SqlMain(){
	fmt.Println("Running sql main : ")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", psqlInfo)
	HandleErr(err)
	err = resetDB(db, dbname)
	HandleErr(err)
	db.Close()

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	HandleErr(err)
	HandleErr(db.Ping())

	db.Close()
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	return err
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	HandleErr(err)
	return createDB(db, name)
}

func HandleErr(err error) {
	if(err != nil){
		panic(err)
	}
}