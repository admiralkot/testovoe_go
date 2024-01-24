package main

// Database: PostgreSQL
//
// Functions to support the interaction with the database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	//"log"
)

// FromJSON decodes a serialized JSON record - User{}
func (p *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// ToJSON encodes a User JSON record
func (p *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// PostgreSQL Connection details
//
// We are using localhost as hostname because both
// the utility and PostgreSQL run on the same machine
var (
	Hostname   = "localhost"
	Port       = 5432
	dbUser     = "postgres"
	dbPassword = "password"
	Database   = "restapi"
)

func ConnectPostgres() *sql.DB {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname, Port, dbUser, dbPassword, Database)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Errorln(err)
		return nil
	}

	return db
}

// DeleteUser is for deleting a user defined by ID
func DeleteUser(ID int) bool {
	db := ConnectPostgres()
	if db == nil {
		log.Errorln("Cannot connect to PostgreSQL!")
		db.Close()
		return false
	}
	defer db.Close()

	// Check is the user ID exists
	t := FindUserID(ID)
	if t.ID == 0 {
		log.Errorln("User", ID, "does not exist.")
		return false
	}

	stmt, err := db.Prepare("DELETE FROM users WHERE ID = $1")
	if err != nil {
		log.Errorln("DeleteUser:", err)
		return false
	}

	_, err = stmt.Exec(ID)
	if err != nil {
		log.Errorln("DeleteUser:", err)
		return false
	}

	return true
}

// InsertUser is for adding a new user to the database
func InsertUser(u User) bool {
	db := ConnectPostgres()
	if db == nil {
		log.Errorln("Cannot connect to PostgreSQL!")
		return false
	}
	defer db.Close()

	if IsUserValid(u) {
		log.Errorln("User", u.Name, "already exists!")
		return false
	}

	stmt, err := db.Prepare("INSERT INTO users(Name, Surname, Patronymic, Age, Gender, Nationality) values($1,$2,$3,$4,$5,$6)")
	if err != nil {
		log.Errorln("Adduser:", err)
		return false
	}

	stmt.Exec(u.Name, u.Surname, u.Patronymic, u.Age, u.Gender, u.Nationality)
	return true
}

// ListAllUsers is for returning all users from the database table
func ListAllUsers() []User {
	db := ConnectPostgres()
	if db == nil {
		log.Errorln("Cannot connect to PostgreSQL!")
		db.Close()
		return []User{}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users \n")
	if err != nil {
		log.Errorln(err)
		return []User{}
	}

	all := []User{}
	var c1 int
	var c2, c3, c4 string
	var c5 int
	var c6, c7 string

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6, c7)
		temp := User{c1, c2, c3, c4, c5, c6, c7}
		all = append(all, temp)
	}

	log.Infoln("All:", all)
	return all
}

// FindUserID is for returning a user record defined by ID
func FindUserID(ID int) User {
	db := ConnectPostgres()
	if db == nil {
		log.Errorln("Cannot connect to PostgreSQL!")
		db.Close()
		return User{}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users where ID = $1\n", ID)
	if err != nil {
		log.Errorln("Query:", err)
		return User{}
	}
	defer rows.Close()

	u := User{}
	var c1 int
	var c2, c3, c4 string
	var c5 int
	var c6, c7 string

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6, &c7)
		if err != nil {
			log.Errorln(err)
			return User{}
		}
		u = User{c1, c2, c3, c4, c5, c6, c7}
		log.Debugln("Found user:", u)
	}
	return u
}

// FindUserUsername is for returning a user record defined by a username
func FindUserName(Name string) User {
	db := ConnectPostgres()
	if db == nil {
		log.Errorln("Cannot connect to PostgreSQL!")
		db.Close()
		return User{}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users where name = $1 \n", Name)
	if err != nil {
		log.Errorln("Fail FindUserName Query:", err)
		return User{}
	}
	defer rows.Close()

	u := User{}
	var c1 int
	var c2, c3, c4 string
	var c5 int
	var c6, c7 string

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6, &c7)
		if err != nil {
			log.Errorln(err)
			return User{}
		}
		u = User{c1, c2, c3, c4, c5, c6, c7}
		log.Debugln("Found user:", u)
	}
	return u
}

func IsUserValid(u User) bool {
	db := ConnectPostgres()
	if db == nil {
		log.Errorln("Cannot connect to PostgreSQL!")
		db.Close()
		return false
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users WHERE name = $1 \n", u.Name)
	if err != nil {
		log.Errorln(err)
		return false
	}

	temp := User{}
	var c1 int
	var c2, c3, c4 string
	var c5 int
	var c6, c7 string

	// If there exist multiple users with the same username,
	// we will get the FIRST ONE only.
	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6, &c7)
		if err != nil {
			log.Errorln(err)
			return false
		}
		temp = User{c1, c2, c3, c4, c5, c6, c7}
	}

	if u.Name == temp.Name && u.Surname == temp.Surname {
		return true
	}
	return false
}

// UpdateUser allows you to update user name
func UpdateUser(u User) bool {
	log.Debugln("Updating user:", u)

	db := ConnectPostgres()
	if db == nil {
		log.Errorln("Cannot connect to PostgreSQL!")
		db.Close()
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE users SET name=$1, surname=$2, patronymic=$3, age=$4, gender=$5, nationality=$6 WHERE ID = $7")
	if err != nil {
		log.Errorln("Adduser:", err)
		return false
	}

	res, err := stmt.Exec(u.Name, u.Surname, u.Patronymic, u.Age, u.Gender, u.Nationality, u.ID)
	if err != nil {
		log.Errorln("UpdateUser failed:", err)
		return false
	}

	affect, err := res.RowsAffected()
	if err != nil {
		log.Errorln("RowsAffected() failed:", err)
		return false
	}
	log.Infoln("Affected:", affect)
	return true
}
