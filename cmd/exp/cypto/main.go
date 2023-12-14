package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func main() {
	f, err := os.OpenFile("passwordGen.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	var in string
	var salt string
	_, err = fmt.Scan(&in, &salt)
	if err != nil {
		panic(err)
	}
	generateFromPassword, err := bcrypt.GenerateFromPassword([]byte(in+salt), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(generateFromPassword))
	fmt.Fprintln(f, string(generateFromPassword))
	fmt.Println(bcrypt.CompareHashAndPassword([]byte("$2a$10$Lt1tVLHyveeF2d1ahaDUjedAvNX80pAUUG9.foGYnYgQINTOOZMRG1"),
		[]byte(string(in+salt))))
}
