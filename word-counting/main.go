package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func producer() chan string {
	f, err := ioutil.ReadFile("words.txt")
	if err != nil {
		fmt.Println("Read failed:", err)
	}
	words := strings.Split(string(f), " ")

	c := make(chan string, 10)
	go func(c chan string) {
		for i := range words {
			c <- words[i]
			fmt.Println("Write:", words[i])
		}
		close(c)
	}(c)
	return c
}

func consumer(stringCh chan string) chan map[string]int {
	c := make(chan map[string]int)
	go func(c chan map[string]int) {
		word_counting := make(map[string]int)
		value, ok := <-stringCh
		for ok {
			time.Sleep(2 * time.Second)
			count, exists := word_counting[value]
			if exists {
				word_counting[value] = count + 1
			} else {
				word_counting[value] = 1
			}
			fmt.Println("Count:", value)
			value, ok = <-stringCh
		}
		c <- word_counting
		close(c)
	}(c)
	return c
}

func mergeChannels(chs ...chan map[string]int) chan map[string]int {
	c := make(chan map[string]int)
	go func(c chan map[string]int) {
		word_counting := make(map[string]int)
		for {
			flag := false
			for i := 0; i < len(chs); i++ {
				value, ok := <-chs[i]
				if ok {
					for k, v := range value {
						count, exists := word_counting[k]
						if exists {
							word_counting[k] = count + v
						} else {
							word_counting[k] = v
						}
					}
					flag = true
				}
			}
			if !flag {
				break
			}
		}
		c <- word_counting
		close(c)
	}(c)
	return c
}

func main() {
	fmt.Println(time.Now())
	c := producer()
	d := mergeChannels(
		consumer(c),
		consumer(c),
		consumer(c),
		consumer(c),
	)
	word_counting := <-d
	for k, v := range word_counting {
		fmt.Println(k+":", v)
	}
	fmt.Println(time.Now())
}
