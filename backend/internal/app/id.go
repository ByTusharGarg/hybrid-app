package app

import (
	"fmt"
	"sync/atomic"
	"time"
)

func newIDGenerator() func(string) string {
	var counter atomic.Uint64
	return func(prefix string) string {
		return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), counter.Add(1))
	}
}
