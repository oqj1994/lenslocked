package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	b := make([]byte, 32)
	_,err:=rand.Read(b)
	if err !=nil {
		panic(err)
	}
	fmt.Println(b)
	fmt.Println(base64.URLEncoding.EncodeToString(b))
}