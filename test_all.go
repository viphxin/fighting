package main

import (
	"fmt"
	"time"
)

func main() {
	st := time.Now()
	//i := 30
	馨 := 1.0 / 30.0 * 1000000000
	c_count := time.Now().Sub(st).Nanoseconds() / int64(馨)
	fmt.Println(c_count)
}
