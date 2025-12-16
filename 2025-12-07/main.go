/*
	simulate tachyon beams moving downward through manifold
	beams split when hitting '^' creating left and right beams
	count total number of splits
*/
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	grid := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, scanner.Text())
	}

	splits := simulateBeams(grid)
	fmt.Printf("Part 1 - Total splits: %d\n", splits)

	timelines := countTimelines(grid)
	fmt.Printf("Part 2 - Total timelines: %d\n", timelines)
}

func simulateBeams(grid []string) int {
	// Find starting position S
	startRow, startCol := -1, -1
	for r, row := range grid {
		for c, ch := range row {
			if ch == 'S' {
				startRow, startCol = r, c
				break
			}
		}
		if startRow != -1 {
			break
		}
	}

	if startRow == -1 {
		return 0
	}

	splits := 0
	hitSplitters := make(map[string]bool)

	// Track active beam positions (use map to merge beams at same position)
	activeBeams := make(map[int]bool)
	activeBeams[startCol] = true

	// Simulate row by row
	for row := startRow; row < len(grid)-1; row++ {
		nextBeams := make(map[int]bool)

		for col := range activeBeams {
			nextRow := row + 1

			// Check bounds
			if col < 0 || col >= len(grid[nextRow]) {
				continue
			}

			// Check what's at next position
			if grid[nextRow][col] == '^' {
				// Hit splitter
				key := fmt.Sprintf("%d,%d", nextRow, col)
				if !hitSplitters[key] {
					hitSplitters[key] = true
					splits++
				}

				// Create two new beams (left and right of splitter)
				if col-1 >= 0 {
					nextBeams[col-1] = true
				}
				if col+1 < len(grid[nextRow]) {
					nextBeams[col+1] = true
				}
			} else {
				// Continue beam downward
				nextBeams[col] = true
			}
		}

		activeBeams = nextBeams
		if len(activeBeams) == 0 {
			break
		}
	}

	return splits
}

func countTimelines(grid []string) int {
	// Find starting position S
	startRow, startCol := -1, -1
	for r, row := range grid {
		for c, ch := range row {
			if ch == 'S' {
				startRow, startCol = r, c
				break
			}
		}
		if startRow != -1 {
			break
		}
	}

	if startRow == -1 {
		return 0
	}

	memo := make(map[string]int)

	var countPaths func(row, col int) int
	countPaths = func(row, col int) int {
		// Exit conditions: reached bottom or out of bounds
		if row >= len(grid) {
			return 1
		}
		if col < 0 || col >= len(grid[row]) {
			return 1
		}

		key := fmt.Sprintf("%d,%d", row, col)
		if val, ok := memo[key]; ok {
			return val
		}

		var result int
		if grid[row][col] == '^' {
			// Splitter: particle takes both paths
			result = countPaths(row+1, col-1) + countPaths(row+1, col+1)
		} else {
			// Empty space or S: continue down
			result = countPaths(row+1, col)
		}

		memo[key] = result
		return result
	}

	return countPaths(startRow, startCol)
}
