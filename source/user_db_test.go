package main

import (
	"testing"
)

func TestUser_UserInit(t *testing.T) {
	err := UserInit()
	if err != nil {
		t.Fail()
	}
}

func TestUser_LoginCheck(t *testing.T) {
	var u User
	u.Email, u.Password = "vkdfg@clc.ua", "sdfsdfvnvd"
	if !u.LoginCheck() {
		t.Fail()
	}
	u.Email, u.Password = "vlad.k@gmail.com", "sertedfsdf"
	if !u.LoginCheck() {
		t.Fail()
	}
}

func TestUser_AuthCheck(t *testing.T) {
	var u User
	u.Email = "vk@clc.com"
	if !u.AuthCheck() {
		t.Fail()
	}
	u.Email = "vlad@gmail.com"
	if !u.AuthCheck() {
		t.Fail()
	}
}

func TestUser_GetUser(t *testing.T) {
	var u User

	u.Email = "vk@clc.com"
	u.GetUser()
	if u.Email == "" {
		t.Fail()
	}

	u.Email = "vlad@gmail.com"
	u.GetUser()
	if u.Email == "" {
		t.Fail()
	}
}

func TestUser_Register(t *testing.T) {
	var userList = []User{
		User{Name: "pass1", Surename: "", Email: "", Password: "", Rights: 0},
		User{Name: "pass1", Surename: "", Email: "", Password: "", Rights: 0},
		User{Name: "pass1", Surename: "", Email: "", Password: "", Rights: 0},
	}

	err := userList[0].Register()
	if err != nil {
		t.Fail()
	}

	err = userList[1].Register()
	if err != nil {
		t.Fail()
	}

	err = userList[2].Register()
	if err != nil {
		t.Fail()
	}
} //
func TestUserField_EditField(t *testing.T) {
	var f = []userField{
		userField{"", "", 0},
		userField{"", "", 0},
		userField{"", "", 0},
	}
	err := f[0].EditField()
	if err != nil {
		t.Fail()
	}
	err = f[1].EditField()
	if err != nil {
		t.Fail()
	}
	err = f[2].EditField()
	if err != nil {
		t.Fail()
	}
} //

func TestUser_IsUsernameAvailable(t *testing.T) {
	var u User
	u.Email = "vk@clc.com"
	if !u.AuthCheck() {
		t.Fail()
	}
	u.Email = "vlad@gmail.com"
	if !u.AuthCheck() {
		t.Fail()
	}
}
