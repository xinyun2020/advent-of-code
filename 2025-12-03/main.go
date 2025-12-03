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
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    part1Total := 0
    part2Total := int64(0)
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        bank := scanner.Text()
        if len(bank) == 0 {
            continue
        }
        part1Total += findMaxJoltage(bank)
        part2Total += findMaxJoltage12(bank)
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading file:", err)
        return
    }

    fmt.Printf("Part 1 - Total output joltage: %d\n", part1Total)
    fmt.Printf("Part 2 - Total output joltage: %d\n", part2Total)
}

func findMaxJoltage(bank string) int {
    max := 0

    for i := 0; i < len(bank); i++ {
        first, _ := strconv.Atoi(string(bank[i]))
        for j := i + 1; j < len(bank); j++ {
            second, _ := strconv.Atoi(string(bank[j]))
            joltage := first*10 + second
            if joltage > max {
                max = joltage
            }
        }
    }

    return max
}

func findMaxJoltage12(bank string) int64 {
    // Select 12 digits to form largest number using greedy stack approach
    toRemove := len(bank) - 12
    result := []byte{}

    for i := 0; i < len(bank); i++ {
        digit := bank[i]
        for len(result) > 0 && toRemove > 0 && result[len(result)-1] < digit {
            result = result[:len(result)-1]
            toRemove--
        }
        result = append(result, digit)
    }

    if toRemove > 0 {
        result = result[:len(result)-toRemove]
    }

    joltage, _ := strconv.ParseInt(string(result), 10, 64)
    return joltage
}
