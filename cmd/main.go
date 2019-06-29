package main

import (
	"fmt"
	"github.com/ShookieShookie/ringslice"
)

func main() {
	fmt.Println("Hello world")
	s := ringslice.NewSlice(10, false)
	fmt.Println(s)

}
