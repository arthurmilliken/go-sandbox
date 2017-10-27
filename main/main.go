package main

import (
	"fmt"
	"sync"
	"time"
)

const workers = 5
const tasks = 100

func main() {
	fmt.Println("main: start")
	start := time.Now()

	cprint := make(chan string)
	done := make(chan bool)
	go printer(cprint, done)

	ctask := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(i, ctask, cprint, &wg)
	}

	for t := 0; t < tasks; t++ {
		ctask <- t
	}
	close(ctask)
	wg.Wait()

	duration := time.Since(start) / time.Millisecond
	cprint <- fmt.Sprintf("main: %d ms", duration)
	close(cprint)
	<- done

	fmt.Println("main: done")
}
