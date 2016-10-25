package main

import (
	"fmt"
	"time"
)

func main() {
	b := 123123
	st := time.Now()
	a := int64(1.0 / 30 * 1000000000)
	c_count := time.Now().Sub(st).Nanoseconds() / a
	fmt.Println(a)
}
