package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type wordCountPair struct {
	Word  string
	Count int
}

func producer() (<-chan string, error) {
	file, err := os.Open("words.txt")
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	c := make(chan string)
	go func() {
		defer close(c)
		for scanner.Scan() {
			c <- scanner.Text()
			fmt.Println("Write:", scanner.Text())
		}
	}()
	return c, nil
}

func consumer(stringCh <-chan string) <-chan wordCountPair {
	c := make(chan wordCountPair)
	go func() {
		defer close(c)
		wordCounting := make(map[string]int)
		for word := range stringCh {
			time.Sleep(2 * time.Second)
			count, exists := wordCounting[word]
			if exists {
				wordCounting[word] = count + 1
			} else {
				wordCounting[word] = 1
			}
		}
		for k, v := range wordCounting {
			c <- wordCountPair{k, v}
		}
	}()
	return c
}

func mergeChannels(chs ...<-chan wordCountPair) <-chan wordCountPair {
	c := make(chan wordCountPair)
	go func() {
		defer close(c)
		var waitGroup sync.WaitGroup
		for i := 0; i < len(chs); i++ {
			waitGroup.Add(1)
			go func(pairCh <-chan wordCountPair) {
				defer waitGroup.Done()
				for pair := range pairCh {
					c <- pair
				}
			}(chs[i])
		}
		waitGroup.Wait()
	}()
	return c
}

func main() {
	fmt.Println(time.Now())
	c, err := producer()
	if err != nil {
		log.Printf("producer failed: %v", err)
	}
	d := mergeChannels(
		consumer(c),
		consumer(c),
		consumer(c),
		consumer(c),
	)
	wordCounting := make(map[string]int)
	for pair := range d {
		count, exists := wordCounting[pair.Word]
		if exists {
			wordCounting[pair.Word] = count + 1
		} else {
			wordCounting[pair.Word] = 1
		}
	}
	for k, v := range wordCounting {
		fmt.Print("{ "+k+": ", v)
		fmt.Println(" }")
	}
	fmt.Println(time.Now())
}
