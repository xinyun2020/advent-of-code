/*
	grid: paper rolls (@) and empty spaces (.)
	part1: count accessible rolls (< 4 neighbors)
	part2: iteratively remove accessible rolls until none remain
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
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	grid := []string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			grid = append(grid, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	part1 := countAccessible(grid)
	fmt.Printf("Part 1 - Accessible rolls: %d\n", part1)

	part2 := countRemovable(grid)
	fmt.Printf("Part 2 - Total removable: %d\n", part2)
}

func countAccessible(grid []string) int {
	return len(filterAccessible(len(grid), len(grid[0]), func(r, c int) bool {
		return grid[r][c] == '@'
	}))
}

func countRemovable(grid []string) int {
	mutableGrid := make([][]rune, len(grid))
	for i, row := range grid {
		mutableGrid[i] = []rune(row)
	}

	totalRemoved := 0
	for {
		accessible := filterAccessible(len(mutableGrid), len(mutableGrid[0]), func(r, c int) bool {
			return mutableGrid[r][c] == '@'
		})
		if len(accessible) == 0 {
			break
		}

		for _, pos := range accessible {
			mutableGrid[pos[0]][pos[1]] = '.'
		}
		totalRemoved += len(accessible)
	}

	return totalRemoved
}

func filterAccessible(rows, cols int, isRoll func(int, int) bool) [][2]int {
	accessible := [][2]int{}
	if rows == 0 || cols == 0 {
		return accessible
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if isRoll(r, c) && isAccessible(r, c, rows, cols, isRoll) {
				accessible = append(accessible, [2]int{r, c})
			}
		}
	}

	return accessible
}

func isAccessible(r, c, rows, cols int, isRoll func(int, int) bool) bool {
	directions := [][2]int{
		{-1, 0}, {-1, 1}, {0, 1}, {1, 1},
		{1, 0}, {1, -1}, {0, -1}, {-1, -1},
	}

	neighbors := 0
	for _, dir := range directions {
		nr, nc := r+dir[0], c+dir[1]
		if nr >= 0 && nr < rows && nc >= 0 && nc < cols && isRoll(nr, nc) {
			neighbors++
		}
	}

	return neighbors < 4
}
