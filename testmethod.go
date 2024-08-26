package main

import "fmt"

func main () {
	type User struct {
		name string 
		surname string
	}

	var John User
	John.name = "John"
	John.surname = "Smith"

	func (u User) Username (John) string {
		return fmt.Printf("Name %v, Surname %v", John.name, John.surname)
	}



}