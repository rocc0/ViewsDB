package main

import (
	"errors"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type (
	//User is used in operations with users, like login, register, etc.
	User struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Surename string `json:"surename"`
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

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) loginCheck() bool {
	var password string

	res := db.QueryRow("SELECT password FROM users WHERE email=$1", u.Email)
	res.Scan(&password)

	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password))

	if err != nil {
		return false
	}
	return true
}

func (u *User) authCheck() bool {
	var privileged int

	res := db.QueryRow("SELECT privileged FROM users WHERE email=$1", u.Email)
	res.Scan(&privileged)

	return privileged == 1
}

func (u *User) getUser() error {

	res := db.QueryRow("SELECT name, surename, id, rights FROM users WHERE email = $1", u.Email)
	err := res.Scan(&u.Name, &u.Surename, &u.ID, &u.Rights)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) register() error {
	if !u.isUsernameAvailable() {
		return errors.New("Користувач з цим ім'ям вже існує")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req, err := db.Prepare("INSERT INTO users (name, surename, email, password) VALUES ($1,$2,$3,$4)")
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
	_, err := stmt.Exec(f.Data, f.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) isUsernameAvailable() bool {
	res, _ := db.Query("SELECT email FROM users WHERE email=$1", u.Email)
	if res == nil {
		return false
	}
	return true
}
