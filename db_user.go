package main

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Surename string `json:"surename"`
	Email    string `json:"email"`
	Password string
	Rights   int
}

func userInit() error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id SERIAL NOT NULL PRIMARY KEY, email VARCHAR(100), " +
		"username VARCHAR(20) NOT NULL, name VARCHAR(100) NOT NULL, password VARCHAR(100) NOT NULL);")
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (u *user) loginCheck() bool {
	var password string

	res := db.QueryRow("SELECT password FROM users WHERE email=?", u.Email)
	res.Scan(&password)

	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password))

	if err != nil {
		return false
	}
	return true
}

func (u *user) authCheck() bool {
	var privileged int

	res := db.QueryRow("SELECT privileged FROM users WHERE email=?", u.Email)
	res.Scan(&privileged)

	return privileged == 1
}

func (u *user) getUser() error {

	res := db.QueryRow("SELECT name, surename, id, rights FROM users WHERE email = ?", u.Email)
	err := res.Scan(&u.Name, &u.Surename, &u.ID, &u.Rights)
	if err != nil {
		return err
	}
	return nil
}

func (u *user) register() error {
	if !u.isUsernameAvailable() {
		return errors.New("Користувач з цим ім'ям вже існує")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req, err := db.Prepare("INSERT INTO users (name, surename, email, password) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = req.Exec(u.Name, u.Surename, u.Email, hashedPassword)

	if err != nil {
		return err
	}

	return nil

}

func (f *userField) editField() error {
	stmt, err := db.Prepare("UPDATE users SET " + f.Field + "=? WHERE id=?;")
	if err != nil {
		return err
	}

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
	_, err = stmt.Exec(f.Data, f.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *user) isUsernameAvailable() bool {
	res, _ := db.Query("SELECT email FROM users WHERE email=?", u.Email)
	if res == nil {
		return false
	}
	return true
}
