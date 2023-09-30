package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_lock(t *testing.T) {

	var count = 0
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Printf("Agent joined the race to aquire the lock\n")
			versionId, err := lock()
			if err != nil {
				fmt.Printf("Agent couldn't aquire the lock, versionId - %s -- %q\n", versionId, err)
			} else {
				mutex.Lock()
				count++
				mutex.Unlock()
			}
		}()
	}
	select {
	case <-time.After(time.Second * 20):
		if count != 1 {
			fmt.Printf("count should always 1  : %d\n", count)
			t.Fail()
		}
	}
}
