package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

type SafeTotal struct {
	val int
	mux sync.Mutex
}

func (total *SafeTotal) Add(count int) {
	total.mux.Lock()
	total.val += count
	total.mux.Unlock()
}

func main() {
	maxRoutine := 5
	total := SafeTotal{val: 0}

	scanner := bufio.NewScanner(os.Stdin)
	wg := new(sync.WaitGroup)
	openChan := make(chan bool, maxRoutine)

	for scanner.Scan() {
		openChan <- true
		wg.Add(1)
		go printResultAndSetTotal(openChan, &total, scanner.Text(), wg)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	wg.Wait()
	close(openChan)
	fmt.Println("Total:", total.val)
}

func printResultAndSetTotal(openChan <-chan bool, total *SafeTotal, url string, wg *sync.WaitGroup) {
	defer wg.Done()

	count := getCountWord(url, "Go")
	fmt.Printf("Count for %s: %d \n", url, count)

	total.Add(count)
	<- openChan
}

func getCountWord(url string, word string) int {
	body := getBodyRequest(url)
	return bytes.Count(body, []byte(word))
}

func getBodyRequest(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}