package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	// Set up the required flags to parse command line arguments.
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	// Open the CSV file specified in the flags and parse it.
	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFileName))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}
	problems := parseLines(lines)

	// Set up a timer with the time limit specified in the flags.
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

	// Loop through each problem and ask the user for an answer. If the timer runs out, break the loop.
problemloop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d. \n", correct, len(problems))
			break problemloop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}

	// Print out the final score.
	fmt.Printf("You scored %d out of %d. \n", correct, len(problems))
}

// Parse each line of the CSV file into a problem struct.
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

// Struct to represent a problem.
type problem struct {
	q string
	a string
}

// Print the provided error message and exit with a status code of 1.
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
