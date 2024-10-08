package main

import (
	"fmt"
	"sync"
)

func main() {
	sum := 0
	var wg sync.WaitGroup
	chans := make(chan int)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			chans <- 0
			sum = sum + 1
			wg.Done()
		}()
		<-chans
	}

	// for i := 0; i < 1000; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		chans <- 0
	// 		sum = sum + 1
	// 		wg.Done()
	// 	}()
	// }

	// for i := 0; i < 1000; i++ {
	// 	<-chans
	// }
	wg.Wait()
	fmt.Println(sum)
}
