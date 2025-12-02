/*
	start with 50, 0-99
	passwaord -> number of times point to 0
	L for minus, R for plus
	if goes below 0, wrap to 99
	if goes above 99, wrap to 0

	use input.txt get final count
*/
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

const dialSize = 100
const startPosition = 50

func main() {
    file, err := os.Open("input.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error opening input.txt: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()

    position := startPosition
    endOnZero := 0
    passZero := 0

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }

        direction := line[0]
        distance, err := strconv.Atoi(line[1:])
        if err != nil {
            fmt.Fprintf(os.Stderr, "Invalid rotation: %s\n", line)
            continue
        }

        var totalPasses int
        if direction == 'L' {
            if position == 0 {
                totalPasses = distance / dialSize
            } else if distance >= position {
                totalPasses = 1 + (distance-position)/dialSize
            }
            position = (position - distance%dialSize + dialSize) % dialSize
        } else {
            totalPasses = (position + distance) / dialSize
            position = (position + distance) % dialSize
        }

        if position == 0 {
            endOnZero++
            if totalPasses > 0 {
                passZero += totalPasses - 1
            }
        } else {
            passZero += totalPasses
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Part 1: %d\n", endOnZero)
    fmt.Printf("Part 2: %d\n", passZero+endOnZero)
}
