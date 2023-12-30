package main

import (
	"fmt"
)

func main() {
	var names []error
	names = append(names, nil)
	for _,n:=range names{
		fmt.Println(n)
	}
}