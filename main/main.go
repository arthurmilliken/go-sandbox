package main

import (
	"fmt"
	"sync"
	"time"
)

const concurrency = 5
const tasks = 100

func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func sleeper(
	id int,
	in <-chan int,
	out chan<- string,
	wg *sync.WaitGroup) {

	out <- fmt.Sprintf("worker %d: start", id)
	for {
		t, ok := <-in
		if ok {
			sleep(100)
			out <- fmt.Sprintf("worker %d: task %d", id, t)
		} else {
			out <- fmt.Sprintf("worker %d: done", id)
			wg.Done()
			return
		}
	}
	out <- fmt.Sprintf("worker %d: done!\n", id)
}

func printer(in <-chan string, done chan<- bool) {
	for {
		str, ok := <-in
		if ok {
			fmt.Println(str)
		} else {
			done <- true
			return
		}
	}
}

func main() {
	fmt.Println("main: start")
	start := time.Now()

	cprint := make(chan string)
	done := make(chan bool)
	go printer(cprint, done)

	ctask := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go sleeper(i, ctask, cprint, &wg)
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
