package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func sleeper(
	id int,
	in <-chan int,
	out chan<- string,
	wg *sync.WaitGroup) {

	log.Printf("sleeper %d: start", id)
	for {
		t, ok := <-in
		if ok {
			sleep(100)
			out <- fmt.Sprintf("sleeper %d: task %2d", id, t)
		} else {
			log.Printf("sleeper %d: done", id)
			wg.Done()
			return
		}
	}
}

func printer(in <-chan string, done chan<- bool) {
	for {
		str, ok := <-in
		if ok {
			log.Println(str)
		} else {
			done <- true
			return
		}
	}
}

func run() {
	log.Println("start")
	start := time.Now()

	cprint := make(chan string)
	done := make(chan bool)
	go printer(cprint, done)

	ctask := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go sleeper(i, ctask, cprint, &wg)
	}

	for t := 0; t < tasks; t++ {
		ctask <- t
	}
	close(ctask)
	wg.Wait()
	close(cprint)
	<-done

	duration := time.Since(start) / time.Millisecond
	log.Printf("%d ms\n", duration)
}
