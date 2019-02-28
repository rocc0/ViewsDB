package main

import (
	"errors"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const userAlreadyExists = "user already exists"

type (
	User struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		SureName string `json:"surname"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Rights   int    `json:"rights"`
	}

	userField struct {
		Field string `json:"field"`
		Data  string `json:"data"`
		ID    int    `json:"id"`
	}
)

func userInit() error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id SERIAL NOT NULL PRIMARY KEY, email VARCHAR(100), " +
		"username VARCHAR(20) NOT NULL, name VARCHAR(100) NOT NULL, password VARCHAR(100) NOT NULL);")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}
	return nil
}

func (u *User) loginCheck() bool {
	var password string

	res := db.QueryRow("SELECT password FROM users WHERE email=$1", u.Email)
	if err := res.Scan(&password); err != nil {
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password)); err != nil {
		return false
	}
	return true
}

func (u *User) authCheck() bool {
	var privileged int

	res := db.QueryRow("SELECT privileged FROM users WHERE email=$1", u.Email)
	if err := res.Scan(&privileged); err != nil {
		return false
	}

	return privileged == 1
}

func (u *User) getUser() error {
	res := db.QueryRow("SELECT name, surename, id, rights FROM users WHERE email = $1", u.Email)
	if err := res.Scan(&u.Name, &u.SureName, &u.ID, &u.Rights); err != nil {
		return err
	}
	return nil
}

func (u *User) register() error {
	if !u.isUsernameAvailable() {
		return errors.New(userAlreadyExists)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req, err := db.Prepare("INSERT INTO users (name, surename, email, password) VALUES ($1,$2,$3,$4)")
	if err != nil {
		return err
	}

	if _, err = req.Exec(u.Name, u.SureName, u.Email, hashedPassword); err != nil {
		return err
	}

	return nil

}

func (f *userField) editField() error {
	stmt, _ := db.Prepare("UPDATE users SET " + f.Field + "=$1 WHERE id=$2;")

	if f.Field == "password" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(f.Data), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(f.Field, hashedPassword, f.ID)
		if err != nil {
			return err
		}
		return nil
	}

	if _, err := stmt.Exec(f.Data, f.ID); err != nil {
		return err
	}
	return nil
}

func (u *User) isUsernameAvailable() bool {

	if res, _ := db.Query("SELECT email FROM users WHERE email=$1", u.Email); res == nil {
		return false
	}
	return true
}
