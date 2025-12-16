/*
Day 12: Christmas Tree Farm - Polyomino Tiling

Fit present shapes (polyominoes) into rectangular regions under trees.
Presents can be rotated and flipped. Solve using backtracking with pruning.
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Coord struct {
	r, c int
}

type Shape []Coord

func parseInput(filename string) ([]Shape, []Region, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var shapes []Shape
	var regions []Region
	var currentShape Shape
	currentShapeRow := 0
	parsingShapes := true

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			if len(currentShape) > 0 {
				shapes = append(shapes, normalizeShape(currentShape))
				currentShape = nil
				currentShapeRow = 0
			}
			continue
		}

		// Check if this is a region line
		if strings.Contains(line, "x") && strings.Contains(line, ":") {
			parsingShapes = false
			if len(currentShape) > 0 {
				shapes = append(shapes, normalizeShape(currentShape))
				currentShape = nil
			}

			parts := strings.Split(line, ":")
			dims := strings.Split(strings.TrimSpace(parts[0]), "x")
			width, _ := strconv.Atoi(dims[0])
			height, _ := strconv.Atoi(dims[1])

			countStrs := strings.Fields(strings.TrimSpace(parts[1]))
			var counts []int
			for _, s := range countStrs {
				count, _ := strconv.Atoi(s)
				counts = append(counts, count)
			}

			regions = append(regions, Region{width: width, height: height, counts: counts})
		} else if parsingShapes {
			// Check if this is a shape index line
			if strings.HasSuffix(strings.TrimSpace(line), ":") {
				if len(currentShape) > 0 {
					shapes = append(shapes, normalizeShape(currentShape))
					currentShape = nil
				}
				currentShapeRow = 0
			} else {
				// Parse shape row
				for c, ch := range line {
					if ch == '#' {
						currentShape = append(currentShape, Coord{r: currentShapeRow, c: c})
					}
				}
				currentShapeRow++
			}
		}
	}

	if len(currentShape) > 0 {
		shapes = append(shapes, normalizeShape(currentShape))
	}

	return shapes, regions, scanner.Err()
}

type Region struct {
	width, height int
	counts        []int
}

func normalizeShape(coords Shape) Shape {
	if len(coords) == 0 {
		return coords
	}

	minR, minC := coords[0].r, coords[0].c
	for _, coord := range coords {
		if coord.r < minR {
			minR = coord.r
		}
		if coord.c < minC {
			minC = coord.c
		}
	}

	normalized := make(Shape, len(coords))
	for i, coord := range coords {
		normalized[i] = Coord{r: coord.r - minR, c: coord.c - minC}
	}
	return normalized
}

func rotate90(shape Shape) Shape {
	rotated := make(Shape, len(shape))
	for i, coord := range shape {
		rotated[i] = Coord{r: coord.c, c: -coord.r}
	}
	return normalizeShape(rotated)
}

func flipHorizontal(shape Shape) Shape {
	flipped := make(Shape, len(shape))
	for i, coord := range shape {
		flipped[i] = Coord{r: coord.r, c: -coord.c}
	}
	return normalizeShape(flipped)
}

func shapeKey(shape Shape) string {
	var sb strings.Builder
	for _, coord := range shape {
		fmt.Fprintf(&sb, "(%d,%d)", coord.r, coord.c)
	}
	return sb.String()
}

func getAllOrientations(shape Shape) []Shape {
	seen := make(map[string]bool)
	var orientations []Shape

	current := normalizeShape(shape)

	for flip := 0; flip < 2; flip++ {
		for rot := 0; rot < 4; rot++ {
			key := shapeKey(current)
			if !seen[key] {
				seen[key] = true
				orientations = append(orientations, current)
			}
			current = rotate90(current)
		}
		current = flipHorizontal(current)
	}

	return orientations
}

func canPlace(grid [][]byte, shape Shape, startR, startC, height, width int) bool {
	for _, coord := range shape {
		r, c := startR+coord.r, startC+coord.c
		if r < 0 || r >= height || c < 0 || c >= width {
			return false
		}
		if grid[r][c] != '.' {
			return false
		}
	}
	return true
}

func placeShape(grid [][]byte, shape Shape, startR, startC int, label byte) {
	for _, coord := range shape {
		r, c := startR+coord.r, startC+coord.c
		grid[r][c] = label
	}
}

func removeShape(grid [][]byte, shape Shape, startR, startC int) {
	for _, coord := range shape {
		r, c := startR+coord.r, startC+coord.c
		grid[r][c] = '.'
	}
}

type Present struct {
	shapeIdx     int
	orientations []Shape
}

func calculateArea(presents []Present, shapes []Shape, start int) int {
	area := 0
	for i := start; i < len(presents); i++ {
		area += len(shapes[presents[i].shapeIdx])
	}
	return area
}

func countEmptySpace(grid [][]byte, height, width int) int {
	count := 0
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			if grid[r][c] == '.' {
				count++
			}
		}
	}
	return count
}

func createGrid(width, height int) [][]byte {
	grid := make([][]byte, height)
	for i := range grid {
		grid[i] = make([]byte, width)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}
	return grid
}

func solveRegion(width, height int, presents []Present, shapes []Shape) bool {
	// Quick feasibility check
	if calculateArea(presents, shapes, 0) > width*height {
		return false
	}

	grid := createGrid(width, height)

	var backtrack func(idx int) bool
	backtrack = func(idx int) bool {
		if idx == len(presents) {
			return true
		}

		// Early termination - check if remaining area is sufficient
		if calculateArea(presents, shapes, idx) > countEmptySpace(grid, height, width) {
			return false
		}

		present := presents[idx]
		label := byte('A' + (idx % 26))

		for _, orientation := range present.orientations {
			for r := 0; r < height; r++ {
				for c := 0; c < width; c++ {
					if canPlace(grid, orientation, r, c, height, width) {
						placeShape(grid, orientation, r, c, label)
						if backtrack(idx + 1) {
							return true
						}
						removeShape(grid, orientation, r, c)
					}
				}
			}
		}

		return false
	}

	return backtrack(0)
}

func solve(filename string) (int, error) {
	shapes, regions, err := parseInput(filename)
	if err != nil {
		return 0, err
	}

	// Precompute all orientations
	allOrientations := make([][]Shape, len(shapes))
	for i, shape := range shapes {
		allOrientations[i] = getAllOrientations(shape)
	}

	count := 0
	for i, region := range regions {
		fmt.Printf("Processing region %d/%d (%dx%d)...\n", i+1, len(regions), region.width, region.height)

		var presents []Present
		for shapeIdx, quantity := range region.counts {
			for j := 0; j < quantity; j++ {
				presents = append(presents, Present{
					shapeIdx:     shapeIdx,
					orientations: allOrientations[shapeIdx],
				})
			}
		}

		if solveRegion(region.width, region.height, presents, shapes) {
			count++
			fmt.Printf("Region %d (%dx%d): FITS\n", i+1, region.width, region.height)
		} else {
			fmt.Printf("Region %d (%dx%d): DOES NOT FIT\n", i+1, region.width, region.height)
		}
	}

	return count, nil
}

func main() {
	filename := "input.txt"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	result, err := solve(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nAnswer: %d regions can fit all their presents\n", result)
}
