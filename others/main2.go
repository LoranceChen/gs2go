package main

import (
	"fmt"
	"unsafe"
)

type user struct {
	userID int
	name   string
	email  string
}

type users []*user

func (users users) String() string {
	s := "["
	for i, user := range users {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", user)
	}
	return s + "]"
}

func Noescape2(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func addUsers(users users) {
	users = append(users, &user{userID: 1, name: "cooluser1", email: "cool.user1@gmail.com"})
	users = append(users, &user{userID: 2, name: "cooluser2", email: "cool.user2@gmail.com"})
	_ = users.String()
	fmt.Printf("users at slice %v \n", users)
}

func main2() {
	var users users
	addUsers(users)
}
