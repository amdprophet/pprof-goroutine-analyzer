package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err.Error())
	os.Exit(1)
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fatal(errors.New("a single argument specifying a path to a pprof goroutine dump is required"))
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fatal(err)
	}
	defer file.Close()

	countsPerMethod := map[string]int{}
	scanner := bufio.NewScanner(file)
	findingNewline := false
	stacktraceLine := 1
	for scanner.Scan() {
		text := scanner.Text()

		if findingNewline {
			if scanner.Text() == "" {
				findingNewline = false
				stacktraceLine = 1
			}
			continue
		}

		if stacktraceLine < 2 {
			stacktraceLine += 1
			continue
		}

		findingNewline = true
		lastOpenParenIndex := strings.LastIndex(text, "(")
		method := text[0:lastOpenParenIndex]
		if count, ok := countsPerMethod[method]; !ok {
			countsPerMethod[method] = 1
		} else {
			countsPerMethod[method] = count + 1
		}
	}

	methodsPerCount := map[int][]string{}
	for method, count := range countsPerMethod {
		methodsPerCount[count] = append(methodsPerCount[count], method)
	}

	counts := []int{}
	for count := range methodsPerCount {
		counts = append(counts, count)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))

	for _, count := range counts {
		for _, method := range methodsPerCount[count] {
			fmt.Printf("%s: %d\n", method, count)
		}
	}
}
