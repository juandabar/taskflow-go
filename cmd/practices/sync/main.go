package main

import (
	"fmt"
	"sync"
	"time"
)

func work1(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(3 * time.Second)
	fmt.Println("work1 finished")
}

func work2(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(3 * time.Second)
	fmt.Println("work2 finished")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go work1(&wg)
	go work2(&wg)

	wg.Wait()

	fmt.Println("Todo ha finalizado...")
}
