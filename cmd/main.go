package main

import (
	"fmt"
	"github.com/ShookieShookie/twincache"
	"time"
)

func main() {
	fmt.Println("Hello world")
	c := twincache.New(10, 5*time.Second)
	for i := 0; i < 10; i++ {
		c.Add(i, i)
	}

	t := time.NewTicker(1 * time.Second)
	go func() {
		for {
			<-t.C
			fmt.Println(c.Expire())
		}
	}()

	time.Sleep(5 * time.Second)
	fmt.Println("\n\n---------------------\nappending again\n")

	for i := 0; i < 20; i++ {
		time.Sleep(1 * time.Second)
		c.Add(i, i)
	}
	time.Sleep(20 * time.Second)
}
