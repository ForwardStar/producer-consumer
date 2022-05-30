package main

import (
	"fmt"
	"producer-consumer/common"
	"time"
)

func main() {
	fmt.Println(time.Now())
	c := common.GenerateInts()
	d := common.MergeChannels(
		common.ConsumeAndDoubleInt(c),
		common.ConsumeAndDoubleInt(c),
		common.ConsumeAndDoubleInt(c),
		common.ConsumeAndDoubleInt(c),
	)
	for i := range d {
		fmt.Println(i)
	}
	fmt.Println(time.Now())
}
