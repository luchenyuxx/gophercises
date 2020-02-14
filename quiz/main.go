package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	csvFile := flag.String("csv", "problems.csv", "the name of the csv file")
	duration := flag.Duration("duration", 30*time.Second, "the duration of the quiz")
	shuffle := flag.Bool("shuffle", false, "shuffle the questions")
	flag.Parse()

	f, err := os.Open(*csvFile)
	if err != nil {
		fmt.Printf("Cannot open file %s\n", *csvFile)
		os.Exit(1)
	}

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if *shuffle {
		rand.Shuffle(len(records), func(i, j int) {
			records[i], records[j] = records[j], records[i]
		})
	}
	scanner := bufio.NewScanner(os.Stdin)
	var totalCount, correctCount int = len(records), 0

	if err != nil {
		fmt.Println("Cannot read csv file")
		os.Exit(1)
	}

	fmt.Println("Hit enter to start, you have", duration)
	scanner.Scan()
	timer := time.NewTimer(*duration)

	for _, s := range records {
		fmt.Println("Question: ", s[0])
		input := make(chan string)
		go func() {
			scanner.Scan()
			input <- strings.TrimSpace(scanner.Text())
		}()
		select {
		case <-timer.C:
			fmt.Println("Time up! Total:", totalCount, ", Correct:", correctCount)
			return
		case i := <-input:
			if i == s[1] {
				correctCount++
				fmt.Println("Right!")
			} else {
				fmt.Println("Wrong!")
			}
		}
	}
}
