package main

import (
	"sync"

	"app"
)

func main() {
	filePath := "input.txt"
	uppercaseFilePath := "uppercase_output.txt"
	reverseFilePath := "reverse_output.txt"

	lines := make(chan string)
	uppercaseLines := make(chan string)
	reverseLines := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)

	go app.ReadLines(filePath, lines)

	go app.ProcessLines(app.UppercaseProcessor{}, lines, uppercaseLines, &wg)
	go app.ProcessLines(app.ReverseProcessor{}, lines, reverseLines, &wg)

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
