package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	tm "github.com/buger/goterm"
)

func displayInitialWords(randomWords []string) {
	tm.Clear()
	tm.Flush()
	fmt.Printf("You have %d seconds to remember the words\n", timeup)
	fmt.Print("======================\n")
	for idx := range randomWords {
		fmt.Print(randomWords[idx] + "\n")
	}
	fmt.Print("======================\n")

	time.Sleep(time.Second * timeup)

	tm.Clear()
	tm.Flush()
}

func shuffle(randomWords []string) {
	rand.Shuffle(
		len(randomWords),
		func(i, j int) { randomWords[i], randomWords[j] = randomWords[j], randomWords[i] })
}

const dictFile = "/usr/share/dict/american-english"
const timeup = 30
const missingWordsNumber = 2

func main() {
	rand.Seed(time.Now().UnixNano())

	info, err := os.Stat(dictFile)
	if os.IsNotExist(err) {
		log.Fatal("Could not find the dict file in ", dictFile)
		os.Exit(1)
	}
	if info.IsDir() {
		log.Fatal("Could not read dict file as it is a directory")
		os.Exit(1)
	}

	// read the file
	content, err := ioutil.ReadFile(dictFile)
	if err != nil {
		log.Fatal("Could not read dict file in ", dictFile)
		os.Exit(1)
	}

	fileData := strings.Split(string(content), "\n")

	data := make([]string, 0)
	for _, line := range fileData {
		if !strings.HasSuffix(line, "'s") {
			data = append(data, line)
		}
	}

	randomWords := make([]string, 0)

	for idx := 0; idx < 10; idx++ {
		randomWord := data[rand.Intn(len(data))]
		randomWords = append(randomWords, randomWord)
	}

	shuffle(randomWords)
	displayInitialWords(randomWords)

	// select some of the words to "remove"
	shuffle(randomWords)
	missingWords := randomWords[len(randomWords)-missingWordsNumber:]
	randomWords = randomWords[:len(randomWords)-missingWordsNumber]

	fmt.Printf("What are the %d missing words?\n", missingWordsNumber)
	fmt.Print("======================\n")
	for idx := range randomWords {
		fmt.Print(randomWords[idx] + "\n")
	}
	fmt.Print("======================\n")

	keepGoing := true
	points := 0

	for keepGoing {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Missing word (BYE for exiting) > ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "BYE" {
			keepGoing = false
		} else {
			valid := false
			for _, missingWord := range missingWords {
				if missingWord == text {
					valid = true
				}
			}

			if valid {
				points = points + 1
				fmt.Printf("Awesome!\n")
			} else {
				fmt.Printf("Nope, try again.\n")
			}

			if points == missingWordsNumber {
				fmt.Print("Got them all!\n")
				keepGoing = false
			}
		}
	}

	if points != missingWordsNumber {
		fmt.Println("Missing words were: ")
		for idx := range missingWords {
			fmt.Printf("-> %s\n", missingWords[idx])
		}
	}
}
