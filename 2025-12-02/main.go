/*
	find invalid product ids in ranges
	invalid ID: sequence of digits repeated twice e.g., 11, 6464, 123123
	sum all invalid ids across all ranges
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func isInvalidPart1(n int) bool {
	s := strconv.Itoa(n)

	if len(s)%2 != 0 {
		return false
	}

	mid := len(s) / 2
	left := s[:mid]
	right := s[mid:]

	return left == right
}

func isInvalidPart2(n int) bool {
	s := strconv.Itoa(n)
	length := len(s)

	for patternLen := 1; patternLen <= length/2; patternLen++ {
		if length%patternLen == 0 {
			pattern := s[:patternLen]
			if strings.Repeat(pattern, length/patternLen) == s {
				return true
			}
		}
	}

	return false
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var input string
	for scanner.Scan() {
		input += scanner.Text() // concat all lines
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	ranges := strings.Split(input, ",")
	sumPart1 := 0
	sumPart2 := 0

	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}

		parts := strings.Split(r, "-")
		if len(parts) != 2 {
			continue
		}

		start, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		end, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		for i := start; i <= end; i++ {
			if isInvalidPart1(i) {
				sumPart1 += i
			}
			if isInvalidPart2(i) {
				sumPart2 += i
			}
		}
	}

	fmt.Printf("Part 1: %d\n", sumPart1)
	fmt.Printf("Part 2: %d\n", sumPart2)
}
