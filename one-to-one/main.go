package main

import (
	"fmt"
	"time"
)

func generateInts() chan int {
	c := make(chan int, 10)
	go func(c chan int) {
		for i := 0; i < 10; i++ {
			c <- i
			fmt.Println("Write", i)
		}
		close(c)
	}(c)
	return c
}

func consumeAndDoubleInt(intCh chan int) chan int {
	c := make(chan int)
	go func(c chan int) {
		value, ok := <-intCh
		for ok {
			time.Sleep(2 * time.Second)
			c <- value * 2
			fmt.Println("Doubled", value)
			value, ok = <-intCh
		}
		close(c)
	}(c)
	return c
}

func main() {
	fmt.Println(time.Now())
	c := generateInts()
	d := consumeAndDoubleInt(c)
	for i := range d {
		fmt.Println(i)
	}
	fmt.Println(time.Now())
}
