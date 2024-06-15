package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Processor interface
type Processor interface {
	Process(data string) string
}

// UppercaseProcessor struct
type UppercaseProcessor struct{}

// ReverseProcessor struct
type ReverseProcessor struct{}

// Process method for UppercaseProcessor
func (u UppercaseProcessor) Process(data string) string {
	return strings.ToUpper(data)
}

// Process method for ReverseProcessor
func (r ReverseProcessor) Process(data string) string {
	runes := []rune(data)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ReadLines reads lines from a file and sends them to a channel
func ReadLines(filePath string, out chan<- string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		close(out)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		out <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	close(out)
}

// ProcessLines processes lines using a given processor and sends the results to a channel
func ProcessLines(processor Processor, in <-chan string, out chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range in {
		out <- processor.Process(line)
	}
}

// WriteLines writes lines from a channel to a file
func WriteLines(in <-chan string, filePath string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for line := range in {
		writer.WriteString(line + "\n")
	}
	writer.Flush()
}
