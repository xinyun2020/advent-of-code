/*
	part 1: check which ingredients fall within ranges
	part 2: merge overlapping ranges and count total fresh IDs
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Range struct {
	start int64
	end   int64
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	ranges := []Range{}
	ingredients := []int64{}
	parsingRanges := true

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			parsingRanges = false
			continue
		}

		if parsingRanges {
			parts := strings.Split(line, "-")
			if len(parts) == 2 {
				start, _ := strconv.ParseInt(parts[0], 10, 64)
				end, _ := strconv.ParseInt(parts[1], 10, 64)
				ranges = append(ranges, Range{start, end})
			}
		} else {
			id, _ := strconv.ParseInt(line, 10, 64)
			ingredients = append(ingredients, id)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	freshCount := 0
	for _, id := range ingredients {
		if isFresh(id, ranges) {
			freshCount++
		}
	}

	fmt.Printf("Part 1 - Fresh ingredients: %d\n", freshCount)

	totalFreshIDs := countTotalFreshIDs(ranges)
	fmt.Printf("Part 2 - Total fresh IDs: %d\n", totalFreshIDs)
}

func isFresh(id int64, ranges []Range) bool {
	for _, r := range ranges {
		if id >= r.start && id <= r.end {
			return true
		}
	}
	return false
}

func countTotalFreshIDs(ranges []Range) int64 {
	if len(ranges) == 0 {
		return 0
	}

	merged := mergeRanges(ranges)

	var count int64
	for _, r := range merged {
		count += r.end - r.start + 1
	}

	return count
}

func mergeRanges(ranges []Range) []Range {
	sorted := make([]Range, len(ranges))
	copy(sorted, ranges)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].start < sorted[j].start
	})

	merged := []Range{sorted[0]}
	for _, current := range sorted[1:] {
		last := &merged[len(merged)-1]

		if current.start <= last.end+1 {
			if current.end > last.end {
				last.end = current.end
			}
		} else {
			merged = append(merged, current)
		}
	}

	return merged
}
