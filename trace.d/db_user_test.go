package main

import (
	"testing"
)

func TestUser_UserInit(t *testing.T) {
	if err := userInit(); err != nil {
		t.Fail()
	}
}

func TestUser_LoginCheck(t *testing.T) {
	u := User{Email: "vkdfg@clc.ua", Password: "sdfsdfvnvd"}
	if !u.loginCheck() {
		t.Fail()
	}
	u = User{Email: "vlad.k@gmail.com", Password: "sertedfsdf"}
	if !u.loginCheck() {
		t.Fail()
	}
}

func TestUser_AuthCheck(t *testing.T) {
	u := User{Email: "vk@clc.com"}
	if !u.authCheck() {
		t.Fail()
	}
	u = User{Email: "vlad@gmail.com"}
	if !u.authCheck() {
		t.Fail()
	}
}

func TestUser_GetUser(t *testing.T) {
	u := User{Email: "vk@clc.com"}
	if err := u.getUser(); err != nil {
		t.Error(err)
	}
	if u.Email == "" {
		t.Fail()
	}
	u = User{Email: "vlad@gmail.com"}
	if err := u.getUser(); err != nil {
		t.Error(err)
	}
	if u.Email == "" {
		t.Fail()
	}
}

func TestUser_Register(t *testing.T) {
	userList := []User{
		{Name: "pass1", SureName: "", Email: "", Password: "", Rights: 0},
		{Name: "pass1", SureName: "", Email: "", Password: "", Rights: 0},
		{Name: "pass1", SureName: "", Email: "", Password: "", Rights: 0},
	}

	if err := userList[0].register(); err != nil {
		t.Fail()
	}

	if err := userList[1].register(); err != nil {
		t.Fail()
	}

	if err := userList[2].register(); err != nil {
		t.Fail()
	}
} //
func TestUserField_EditField(t *testing.T) {
	fields := []userField{
		{"", "", 0},
		{"", "", 0},
		{"", "", 0},
	}

	if err := fields[0].editField(); err != nil {
		t.Fail()
	}

	if err := fields[1].editField(); err != nil {
		t.Fail()
	}

	if err := fields[2].editField(); err != nil {
		t.Fail()
	}
} //

func TestUser_IsUsernameAvailable(t *testing.T) {
	u := User{Email: "vk@clc.com"}
	if !u.authCheck() {
		t.Fail()
	}
	u = User{Email: "vlad@gmail.com"}
	if !u.authCheck() {
		t.Fail()
	}
}
