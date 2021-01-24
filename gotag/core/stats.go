package core

import (
	"fmt"
	"time"
)

func TimeStart() (start int64) {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func TimeStop(start int64) {
	stop := time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	timeTaken := stop - start
	fmt.Printf("%dms\n", timeTaken)
}
