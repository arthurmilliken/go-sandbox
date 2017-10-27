package main

import (
	"log"
)

const workers = 5
const tasks = 100

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile | log.Lmicroseconds)
	walk()
}
