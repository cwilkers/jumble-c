package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

type Data struct {
	Answer []string
}

// Vars that need to be global
var wordHash map[string][]string

// Hash function that alphabetizes letters in a word
func hash(s string) string {
	a := strings.Split(s, "")
	sort.Strings(a)
	return strings.Join(a, "")
}

// Parse dictionary into hash of slices
func readWords(wordFile string) map[string][]string {
	var wordHash = make(map[string][]string)
	wordList, err := os.Open(wordFile)
	if err != nil {
		log.Fatal(err)
	}
	defer wordList.Close()
	scanner := bufio.NewScanner(wordList)
	for scanner.Scan() {
		myWord := strings.Trim(scanner.Text(), " ")
		myHash := hash(myWord)
		wordHash[myHash] = append(wordHash[myHash], myWord)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return wordHash
}

// Subdivide guesses to come up with more answers
func divide(q string) []string {
    l := len(q)
	var answer []string
	if l > 2 {
        a := strings.Split(q, "")
		append(answer, wordHash[strings.Join(a[1:l-1], "")])
        for i:= 1; i < l-2; i++ {
		    append(answer, wordHash[strings.Join(a[0:i]+a[i+1,l], "")])
		}
		append(answer, wordHash[strings.Join(a[0:l-2], "")])
	}
	return answer
}

// handle / and run queries
func formHandler(w http.ResponseWriter, r *http.Request) {
    query := r.FormValue("search")
	answer:= wordHash[hash(query)]
	d := &Data{Answer: answer}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, d)
}

// Set up a webserver
func main() {
	// Check environment variables
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	address := fmt.Sprintf("%s:%s", "0.0.0.0", httpPort)
	wordFile := os.Getenv("WORDFILE")
	if wordFile == "" {
		wordFile = "words"
	}
	wordHash = readWords(wordFile)
	fmt.Printf("Read %v words\n", len(wordHash))
    http.HandleFunc("/", formHandler)
	log.Fatal(http.ListenAndServe(address, nil))
}