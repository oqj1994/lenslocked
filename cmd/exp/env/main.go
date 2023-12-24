package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err !=nil{
		panic(err)
	}

	name:=os.Getenv("name")
	fmt.Println(name)
}