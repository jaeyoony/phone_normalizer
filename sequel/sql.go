package sequel

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/jaeyoony/phone_normalizer/normalizer"
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
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	HandleErr(err)
	defer db.Close()
	HandleErr(db.Ping())

	// create the numbers table 
	HandleErr(createNumbersTable(db))
	var id int
	id, err = insertPhoneNumber(db, "1234567890")
	HandleErr(err)

	// var number string
	// number, err = getPhone(db, id)
	// HandleErr(err)

	// remove duplicates of the default number
	err = normalize(db, id)
	HandleErr(err)

	// insert new whacky number 
	id, err = insertPhoneNumber(db, "987as-6534-3!21f0")
	HandleErr(err)

	// get all numbers, for sanity check 
	phones, err := getAllNumbers(db)
	HandleErr(err)
	fmt.Println(phones)

	// get all ids to sanitize 
	all_ids, err := getAllIds(db)
	fmt.Println("all ids : ", all_ids)
	HandleErr(err)
	for _, id := range all_ids {
		err = removeIfIsDuplicate(db, id)
		HandleErr(err)
		err = normalize(db, id)
		HandleErr(err)
	}

	// get all numbers, for sanity check 
	phones, err = getAllNumbers(db)
	HandleErr(err)
	fmt.Println(phones)


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

func createNumbersTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL, 
			value VARCHAR(255)
	)`

	_, err := db.Exec(statement)
	return err
}

func insertPhoneNumber(db *sql.DB, phone string) (int, error) {
	// return the id of the entry in the phone_numbers table
	//	and error if there is one 

	// the weird dollar sign notation here is part of the postgres driver, 
	//	and is used to prevent injection attacks
	//	basically, when we call the Exec statement onto the db,
	//	we pass in the arguments in order, and the $n is used to reference
	//	the nth parameter we pass in after our statement. 
	//	in this case, it'd be the phone number bc that is the 1st param
	//	after the statement 
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	
	var id int
	// we have to use QueryRow() here because for some reason postgres doesn't 
	//	return the id of the newly inserted row
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	err := db.QueryRow("SELECT value FROM phone_numbers WHERE id=$1", id).Scan(&number)
	if err != nil {
		return "", err
	}
	return number, nil
}

func getAllNumbers(db *sql.DB) ([]string, error) {
	var all_numbers []string

	// the db.Query() fcn returns a number of rows
	//	ref : https://pkg.go.dev/github.com/lib/pq
	// to go through all the returned rows, we have to use the 
	//	rows.Next() to loop thru the results. 
	//	from there, we treat each individual row like we did before with our
	//	single row queries, where we scan for values we want and assign to vars. 
	rows, err := db.Query("SELECT value FROM phone_numbers")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var number string
		if err := rows.Scan(&number); err != nil {
			return nil, err
		}
		all_numbers = append(all_numbers, number)
	}

	// remember, after we go thru all rows, we have to check if there's
	//	an error at the end using rows.Err()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return all_numbers, nil
}

func getAllIds(db *sql.DB) ([]int, error) {
	var all_ids []int
	rows, err := db.Query("SELECT id FROM phone_numbers")
	if err != nil  {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		all_ids = append(all_ids, id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return all_ids, nil
}

func getIdByNumber(db *sql.DB, number string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM phone_numbers WHERE value=$1", number).Scan(&id)
	if err == sql.ErrNoRows {
		return -1, nil
	}
	if err != nil {
		return -1, err
	}
	return id, nil
}

func normalize(db *sql.DB, num_id int) error {
	// uses the normalize fcn we created in teh normalizer package
	//	to normalize the phone number
	var num string
	err := db.QueryRow("SELECT value FROM phone_numbers WHERE id=$1", num_id).Scan(&num)
	HandleErr(err)
	new_num := normalizer.NormalizeNumber(num)

	if new_num != num {
		// if this number already exists, remove from database
		temp_id, err := getIdByNumber(db, new_num)
		HandleErr(err)
		fmt.Println("Temp id : ", temp_id)

		if(temp_id != -1 && temp_id != num_id) {
			db.Exec("DELETE FROM phone_numbers WHERE id=$1", num_id)
			fmt.Println("deleted row because normalized version already exists")
		} else {
			db.Exec("UPDATE phone_numbers SET value=$1 WHERE id=$2", new_num, num_id)
			fmt.Println("Updated to normalized number")
		}
	}
	return nil
}

func removeIfIsDuplicate(db *sql.DB, num_id int) error {
	var num string
	err := db.QueryRow("SELECT value FROM phone_numbers WHERE id=$1", num_id).Scan(&num)
	// if(err == sql.ErrNoRows){
	// 	return nil
	// }

	if(err!=nil){
		return err
	}

	temp_id, err := getIdByNumber(db, num)
	if err != nil {
		return err
	}

	if(temp_id != -1 && temp_id != num_id){
		db.Exec("DELETE FROM phone_numbers WHERE id=$1", num_id)
		fmt.Println("removed duplicate")
	}
	return nil
}