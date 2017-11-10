package main

import (
	"fmt"
	"time"
)

func timeStart() (start int64) {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func timeStop(start int64) {
	stop := time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	timeTaken := stop - start
	fmt.Printf("%dms\n", timeTaken)
}
