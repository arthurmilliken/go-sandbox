package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
    "crypto/sha256"
)

const defaultHash = "----------------------------------------------"

type Node struct {
	path string
	info os.FileInfo
	hash string
}

func (n *Node) String() string {
	return fmt.Sprintf("%s %s %s (%d bytes)",
		n.hash,
		n.info.Mode(),
		n.path,
		n.info.Size())
}

func ipfsAddress(path string) string {
	cmd := exec.Command("ipfs", "add", "-nrQ", path)
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("%s: %v", path, err)
	}
	return strings.TrimSpace(string(out[:]))
}

func sha(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		content = []byte(path)
	}
	sum := sha256.Sum256(content)
    return fmt.Sprintf("%x", sum)
}

func worker(in <-chan *Node, out chan<- *Node, wg *sync.WaitGroup) {
	for {
		node, ok := <- in
		if ok {
			node.hash = sha(node.path)
			out <- node
		} else {
			wg.Done()
			return
		}
	}
}

func receiver(in <-chan *Node, done chan<- bool) {
	for {
		node, ok := <- in
		if ok {
			fmt.Println(node)
			// log.Print(node)
		} else {
			done <- true
			return
		}
	}

}

func walk() {
	log.Println("start")
	start := time.Now()
	conf := LoadConfig()
	root := conf.DefaultPath
	log.Printf("%+v", conf)

	hashed := make(chan *Node)
	done := make(chan bool)
	go receiver(hashed, done)

	nodes := make(chan *Node)
	var wg sync.WaitGroup
	for i := 0; i < conf.Workers; i++ {
		wg.Add(1)
		go worker(nodes, hashed, &wg)
	}

	step := func(path string, info os.FileInfo, err error) error {
		nodes <- &Node{path: path, info: info}
		// log.Println(&Node{path: path, info: info, hash: ipfsAddress(path)})
		return err
	}
	err := filepath.Walk(root, step)
	if err != nil {
		log.Fatal(err)
	}
	close(nodes)
	wg.Wait()
	close(hashed)
	<-done
	log.Printf("%d ms\n", time.Since(start) / time.Millisecond)
}
