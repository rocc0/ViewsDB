package main

import (
	"errors"
	"strings"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Username string
	Password string
	Name string
	Email string
}

func userInit() {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id SERIAL NOT NULL PRIMARY KEY, email VARCHAR(100), " +
		"username VARCHAR(20) NOT NULL, name VARCHAR(100) NOT NULL, password VARCHAR(100) NOT NULL);")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec()

	if err != nil {
		fmt.Print(err.Error())
	}
}


func isUserValid(username, password string) bool {
	var databaseUsername  string
	var databasePassword  string
	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func registerNewUser(username, password, name, email string) (*user, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !isUsernameAvailable(username) {
		return nil, errors.New("The username isn't available")
	}
	u := user{Username: username, Password: password, Email: email, Name: name}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("INSERT INTO users(username, password, name, email) VALUES(?, ?, ?, ?)", username, hashedPassword, name, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func isUsernameAvailable(username string) bool {
	_, err := db.Query("SELECT username FROM users WHERE username=?", username)
	if err != nil {
		return false
	}
	return true
}