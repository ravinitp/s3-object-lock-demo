package main

import (
	"fmt"
	"sync"
	"testing"
)

func Test_lock(t *testing.T) {
	bucket := "lock-demo"
	path := "terraform.tfstate"
	var LockVersionId = ""
	mutex := sync.Mutex{}
	var wg sync.WaitGroup
	var count = 0

	lockRecieved := make(chan bool)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("Agent joined the race to aquire the lock\n")
			versionId, err := lock(bucket, path)
			if err != nil {
				fmt.Printf("Agent couldn't aquire the lock, versionId - %s -- %q\n", versionId, err)
				lockRecieved <- false
			} else {
				mutex.Lock()
				LockVersionId = versionId
				count++
				lockRecieved <- true
				mutex.Unlock()
			}
		}()
	}

	go func() {
		wg.Wait()
		close(lockRecieved)
		unlock(LockVersionId, bucket, path)
	}()

	var total int64
	for lr := range lockRecieved {
		if lr {
			total++
		}
	}
	if total != 1 && count != 1 {
		fmt.Printf("total: %d, count: %d\n", total, count)
		t.Fail()
	}
}
