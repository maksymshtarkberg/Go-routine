package main

import (
	"sync"

	"github.com/maksymshtarkberg/Go-routine/app"
)

func main() {
	filePath := "input.txt"
	uppercaseFilePath := "uppercase_output.txt"
	reverseFilePath := "reverse_output.txt"

	lines := make(chan string)
	uppercaseLines := make(chan string)
	reverseLines := make(chan string)

	var wg sync.WaitGroup
	wg.Add(3)

	// Intermediate channels for separate processing
	uppercaseInput := make(chan string)
	reverseInput := make(chan string)

	go app.ReadLines(filePath, lines)

	// Intermediate goroutine to split input to two separate channels
	go func() {
		defer close(uppercaseInput)
		defer close(reverseInput)
		for line := range lines {
			uppercaseInput <- line
			reverseInput <- line
		}
		wg.Done()
	}()

	go app.ProcessLines(app.UppercaseProcessor{}, uppercaseInput, uppercaseLines, &wg)
	go app.ProcessLines(app.ReverseProcessor{}, reverseInput, reverseLines, &wg)

	go func() {
		wg.Wait()
		close(uppercaseLines)
		close(reverseLines)
	}()

	var writeWg sync.WaitGroup
	writeWg.Add(2)
	go app.WriteLines(uppercaseLines, uppercaseFilePath, &writeWg)
	go app.WriteLines(reverseLines, reverseFilePath, &writeWg)

	writeWg.Wait()
}
