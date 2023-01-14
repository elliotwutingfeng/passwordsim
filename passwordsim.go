// Package passwordsim lets you search for passwords similar to
// your specified password in any passwords dataset.
// The similarity metric used is the Damerau-Levenshtein distance.
package passwordsim

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"

	"github.com/fatih/color"
	tdl "github.com/lmas/Damerau-Levenshtein"
)

type results []result

func (u results) Len() int {
	return len(u)
}
func (u results) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}
func (u results) Less(i, j int) bool {
	if u[i].score < u[j].score {
		return true
	}
	if u[i].score == u[j].score {
		return u[i].password < u[j].password
	}
	return false
}

type result struct {
	password string
	score    float64
}

// bufferSize is buffer capacity of message channel
//
// CPU-bound receivers are slower at calculating distances
// than I/O sender at sending out passwords, so bufferSize
// has to be high enough.
const bufferSize = 512

// estimatedMaxPasswordLength sets an optimal initial size for package tdl's distance calculator
// to reduce chances of having to resize (i.e. most passwords do not exceed 64 characters).
const estimatedMaxPasswordLength = 64

// CheckPasswords scans passwords dataset for passwords similar to passwordToCheck,
// where the normalised Damerau-Levenshtein distance is no higher than threshold, and writes the passwords and their
// respective normalised Damerau-Levenshtein distance scores to file output.
func CheckPasswords(passwords string, output string, passwordToCheck string, threshold float64) {
	var numReceivers = runtime.NumCPU()
	var wg sync.WaitGroup
	var similarPasswords sync.Map

	if threshold < 0 {
		threshold = 0
	}
	if threshold > 1 {
		threshold = 1
	}

	message := make(chan string, bufferSize)
	go send(message, passwords)

	for i := 0; i < numReceivers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t := tdl.New(estimatedMaxPasswordLength)
			for password := range message {
				score := t.Distance(password, passwordToCheck)
				normalisedScore := float64(score) / float64(math.Max(float64(len(password)), float64(len(passwordToCheck))))
				if normalisedScore <= threshold {
					similarPasswords.Store(password, normalisedScore)
				}
			}
		}()
	}
	wg.Wait()
	writeFile, err := os.OpenFile(output, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		color.HiRed("Error: %s", err)
		os.Exit(1)
	}
	defer writeFile.Close()

	res := results{}
	similarPasswords.Range(func(k, v interface{}) bool {
		res = append(res, result{password: k.(string), score: v.(float64)})
		return true
	})
	sort.Sort(res)

	save(res, writeFile)

	var passwordCount string
	var clr *color.Color
	if len(res) == 0 {
		passwordCount = "No similar passwords found"
		clr = color.New(color.Bold, color.FgHiGreen)
	} else if len(res) == 1 {
		passwordCount = "1 similar password found"
		clr = color.New(color.Bold, color.FgRed)
	} else {
		passwordCount = strconv.Itoa(len(res)) + " similar passwords found"
		clr = color.New(color.Bold, color.FgRed)
	}
	clr.Print(passwordCount + " in '" + passwords +
		"'. Threshold: " + strconv.FormatFloat(threshold, 'f', 2, 64) + "\n")
	color.White("Results saved to file '%s'.\n", output)
}

// save saves passwords and their normalised Damerau-Levenshtein distance scores to file opened at writeFile
func save(res results, writeFile *os.File) {
	for _, r := range res {
		fmt.Fprintln(writeFile, r.password, r.score)
	}
}

// send streams passwords from passwords dataset to channel ch
func send(ch chan<- string, passwords string) {
	readFile, err := os.Open(passwords)
	if err != nil {
		color.HiRed("Error: %s", err)
		os.Exit(1)
	}
	defer readFile.Close()
	reader := bufio.NewScanner(readFile)
	color.HiYellow("Searching for similar passwords...")

	for reader.Scan() {
		ch <- reader.Text()
	}
	close(ch)
}
