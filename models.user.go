package main

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Surename string `json:"surename"`
	Email string `json:"email"`
	Rights int
}

func userInit() {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id SERIAL NOT NULL PRIMARY KEY, email VARCHAR(100), " +
		"username VARCHAR(20) NOT NULL, name VARCHAR(100) NOT NULL, password VARCHAR(100) NOT NULL);")
	check(err)

	_, err = stmt.Exec()
	check(err)
}


func loginCheck(email, password string) bool{
	var (
		e_mail, passw string
	)
	res := db.QueryRow("SELECT email, password FROM users WHERE email=?", email)
	res.Scan(&e_mail, &passw)
	err := bcrypt.CompareHashAndPassword([]byte(passw), []byte(password))

	if err != nil {
		return false
	}
	return true
}

func authCheck(email string) bool {
	var (
		privileged int
	)
	res := db.QueryRow("SELECT privileged FROM users WHERE email=?", email)
	res.Scan(&privileged)

	return privileged == 1
}

func getUser(email string) (*User, error) {

	var (
		userData User
		name, surename, e_mail string
		id,rights int
	)

	res := db.QueryRow("SELECT name, surename, email, id, rights FROM users WHERE email = ?", email)
	err := res.Scan(&name, &surename, &e_mail, &id, &rights)
	if err != nil {
		return &userData, err
	}

	userData = User{id, name, surename, e_mail, rights}

	return &userData, nil
}

func postUser(name, surename, email, password string) (string, error) {
	if !isUsernameAvailable(email) {

		return name, errors.New("Користувач з цим ім'ям вже існує")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	req, err := db.Prepare("INSERT INTO users (name, surename, email, password) VALUES (?,?,?,?)")
	check(err)
	_, err = req.Exec(name, surename, email, hashedPassword)
	check(err)

	return name, nil

}

func postChangeField(field, data string, id int) error {
	stmt, err := db.Prepare("UPDATE users SET "+ field + "=? WHERE id=?;")
	check(err)

	if field == "password" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
		_, err = stmt.Exec(field, hashedPassword, id)
		check(err)
		return nil
	} else {
		_, err = stmt.Exec(data, id)
		check(err)
		return nil
	}
}


func isUsernameAvailable(email string) bool {
	res, _ := db.Query("SELECT email FROM users WHERE email=?", email)
	if res == nil {

		return false
	}
	return true
}