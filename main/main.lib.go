package main

import (
	"fmt"
	"sync"
	"time"
)

func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func worker(
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
