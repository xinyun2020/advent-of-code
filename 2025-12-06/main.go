/*
	parse vertical math worksheet
	problems separated by full column of spaces
	part 2: read columns right-to-left, digits top-to-bottom
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	rows := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rows = append(rows, scanner.Text())
	}

	part1 := parseWorksheetPart1(rows)
	fmt.Printf("Part 1: %d\n", part1)

	part2 := parseWorksheetPart2(rows)
	fmt.Printf("Part 2: %d\n", part2)
}

func parseWorksheetPart1(rows []string) int64 {
	if len(rows) < 5 {
		return 0
	}

	maxLen := 0
	for _, row := range rows {
		if len(row) > maxLen {
			maxLen = len(row)
		}
	}

	separators := []int{-1}
	for col := 0; col < maxLen; col++ {
		allSpaces := true
		for _, row := range rows {
			if col < len(row) && row[col] != ' ' {
				allSpaces = false
				break
			}
		}
		if allSpaces {
			separators = append(separators, col)
		}
	}
	separators = append(separators, maxLen)

	var grandTotal int64
	for i := 0; i < len(separators)-1; i++ {
		start := separators[i] + 1
		end := separators[i+1]

		if start >= end {
			continue
		}

		problemRows := make([]string, len(rows))
		for r, row := range rows {
			if start < len(row) {
				if end <= len(row) {
					problemRows[r] = row[start:end]
				} else {
					problemRows[r] = row[start:]
				}
			}
		}

		result := solveProblemPart1(problemRows)
		grandTotal += result
	}

	return grandTotal
}

func solveProblemPart1(rows []string) int64 {
	if len(rows) < 5 {
		return 0
	}

	operator := byte(0)
	operatorRow := rows[4]
	for _, ch := range operatorRow {
		if ch == '+' || ch == '*' {
			operator = byte(ch)
			break
		}
	}

	if operator == 0 {
		return 0
	}

	numbers := []int64{}
	for i := 0; i < 4; i++ {
		numStr := ""
		for _, ch := range rows[i] {
			if ch >= '0' && ch <= '9' {
				numStr += string(ch)
			}
		}
		if numStr != "" {
			num, _ := strconv.ParseInt(numStr, 10, 64)
			numbers = append(numbers, num)
		}
	}

	if len(numbers) == 0 {
		return 0
	}

	result := numbers[0]
	for i := 1; i < len(numbers); i++ {
		if operator == '+' {
			result += numbers[i]
		} else {
			result *= numbers[i]
		}
	}

	return result
}

func parseWorksheetPart2(rows []string) int64 {
	if len(rows) < 5 {
		return 0
	}

	maxLen := 0
	for _, row := range rows {
		if len(row) > maxLen {
			maxLen = len(row)
		}
	}

	separators := []int{-1}
	for col := 0; col < maxLen; col++ {
		allSpaces := true
		for _, row := range rows {
			if col < len(row) && row[col] != ' ' {
				allSpaces = false
				break
			}
		}
		if allSpaces {
			separators = append(separators, col)
		}
	}
	separators = append(separators, maxLen)

	var grandTotal int64
	for i := 0; i < len(separators)-1; i++ {
		start := separators[i] + 1
		end := separators[i+1]

		if start >= end {
			continue
		}

		problemRows := make([]string, len(rows))
		for r, row := range rows {
			if start < len(row) {
				if end <= len(row) {
					problemRows[r] = row[start:end]
				} else {
					problemRows[r] = row[start:]
				}
			}
		}

		result := solveProblemPart2(problemRows)
		grandTotal += result
	}

	return grandTotal
}

func solveProblemPart2(rows []string) int64 {
	if len(rows) < 5 {
		return 0
	}

	// Find operator
	operator := byte(0)
	operatorRow := rows[4]
	for _, ch := range operatorRow {
		if ch == '+' || ch == '*' {
			operator = byte(ch)
			break
		}
	}

	if operator == 0 {
		return 0
	}

	// Read columns right-to-left
	numbers := []int64{}
	width := len(rows[0])
	for col := width - 1; col >= 0; col-- {
		// Read digits top-to-bottom in this column
		numStr := ""
		for row := 0; row < 4; row++ {
			if col < len(rows[row]) {
				ch := rows[row][col]
				if ch >= '0' && ch <= '9' {
					numStr += string(ch)
				}
			}
		}

		if numStr != "" {
			num, _ := strconv.ParseInt(numStr, 10, 64)
			numbers = append(numbers, num)
		}
	}

	if len(numbers) == 0 {
		return 0
	}

	result := numbers[0]
	for i := 1; i < len(numbers); i++ {
		if operator == '+' {
			result += numbers[i]
		} else {
			result *= numbers[i]
		}
	}

	return result
}
