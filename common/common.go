package common

import (
	"fmt"
	"sync"
	"time"
)

func GenerateInts() <-chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for i := 0; i < 10; i++ {
			c <- i
			fmt.Println("Write", i)
		}
	}()
	return c
}

func ConsumeAndDoubleInt(intCh <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for i := range intCh {
			time.Sleep(2 * time.Second)
			c <- i * 2
			fmt.Println("Doubled", i)
		}
	}()
	return c
}

func MergeChannels(chs ...<-chan int) <-chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		var waitGroup sync.WaitGroup
		for i := 0; i < len(chs); i++ {
			waitGroup.Add(1)
			go func(intCh <-chan int) {
				defer waitGroup.Done()
				for value := range intCh {
					c <- value
				}
			}(chs[i])
		}
		waitGroup.Wait()
	}()
	return c
}
