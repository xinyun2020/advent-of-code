/*
Day 10: Factory

Part 1: GF(2) Gaussian elimination (binary field)
Part 2: Integer Gaussian elimination with bounded free variable search
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	lightsRe  = regexp.MustCompile(`\[(.*?)\]`)
	buttonsRe = regexp.MustCompile(`\(([0-9,]+)\)`)
	joltageRe = regexp.MustCompile(`\{([0-9,]+)\}`)
)

type Machine struct {
	lights  []int
	joltage []int
	buttons [][]int
}

func main() {
	machines, err := parseInput("input.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var part1, part2 int
	for i, m := range machines {
		part1 += solvePart1(m)
		part2 += solvePart2(m, i+1)
	}

	fmt.Printf("\nPart 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

func parseInput(filename string) ([]Machine, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var machines []Machine
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			machines = append(machines, Machine{
				lights:  parseLights(line),
				joltage: parseInts(joltageRe, line),
				buttons: parseAllInts(buttonsRe, line),
			})
		}
	}
	return machines, scanner.Err()
}

func parseLights(line string) []int {
	match := lightsRe.FindStringSubmatch(line)
	if len(match) < 2 {
		return nil
	}
	lights := make([]int, len(match[1]))
	for i, c := range match[1] {
		if c == '#' {
			lights[i] = 1
		}
	}
	return lights
}

func parseInts(re *regexp.Regexp, line string) []int {
	match := re.FindStringSubmatch(line)
	if len(match) < 2 {
		return nil
	}
	return splitInts(match[1])
}

func parseAllInts(re *regexp.Regexp, line string) [][]int {
	matches := re.FindAllStringSubmatch(line, -1)
	result := make([][]int, 0, len(matches))
	for _, m := range matches {
		if len(m) >= 2 {
			result = append(result, splitInts(m[1]))
		}
	}
	return result
}

func splitInts(s string) []int {
	var result []int
	for _, part := range strings.Split(s, ",") {
		if n, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			result = append(result, n)
		}
	}
	return result
}

// Part 1: GF(2) - binary field where 1+1=0
func solvePart1(m Machine) int {
	matrix := buildMatrix(m.buttons, len(m.lights), false)
	return solveGF2(matrix, m.lights)
}

// Part 2: Integer linear system
func solvePart2(m Machine, machineNum int) int {
	n, nButtons := len(m.joltage), len(m.buttons)
	coeff := buildMatrix(m.buttons, n, true)
	aug := augment(coeff, m.joltage)
	pivots := rref(aug, nButtons)
	freeVars := findFreeVars(pivots, nButtons)

	result := searchMinSolution(aug, pivots, freeVars, nButtons, slices.Max(m.joltage))

	if result == -1 {
		fmt.Printf("Machine %d: 0 (no solution)\n", machineNum)
		return 0
	}
	fmt.Printf("Machine %d: %d\n", machineNum, result)
	return result
}

// buildMatrix creates coefficient matrix from button mappings.
// If additive=true, counts occurrences; if false, sets to 1 (binary).
func buildMatrix(buttons [][]int, rows int, additive bool) [][]int {
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, len(buttons))
	}
	for j, button := range buttons {
		for _, idx := range button {
			if idx < rows {
				if additive {
					matrix[idx][j]++
				} else {
					matrix[idx][j] = 1
				}
			}
		}
	}
	return matrix
}

func augment(matrix [][]int, target []int) [][]int {
	aug := make([][]int, len(matrix))
	for i, row := range matrix {
		aug[i] = make([]int, len(row)+1)
		copy(aug[i], row)
		aug[i][len(row)] = target[i]
	}
	return aug
}

// rref performs row reduction (RREF) over integers
func rref(aug [][]int, nCols int) []int {
	var pivots []int
	row := 0

	for col := 0; col < nCols && row < len(aug); col++ {
		pivotRow := findPivot(aug, row, col, false)
		if pivotRow == -1 {
			continue
		}

		aug[row], aug[pivotRow] = aug[pivotRow], aug[row]
		pivots = append(pivots, col)

		for r := range aug {
			if r != row && aug[r][col] != 0 {
				f1, f2 := aug[r][col], aug[row][col]
				for c := 0; c <= nCols; c++ {
					aug[r][c] = aug[r][c]*f2 - aug[row][c]*f1
				}
			}
		}
		row++
	}
	return pivots
}

func searchMinSolution(aug [][]int, pivots, freeVars []int, nButtons, maxVal int) int {
	minPresses := -1
	assignment := make([]int, len(freeVars))

	var search func(int)
	search = func(idx int) {
		if idx == len(freeVars) {
			solution := make([]int, nButtons)
			for i, v := range freeVars {
				solution[v] = assignment[i]
			}
			if backSubstitute(aug, pivots, solution, nButtons) {
				if sum := sum(solution); minPresses == -1 || sum < minPresses {
					minPresses = sum
				}
			}
			return
		}

		bound := maxVal
		if minPresses != -1 {
			bound = min(maxVal, minPresses-sum(assignment[:idx])-1)
		}

		for val := 0; val <= bound; val++ {
			assignment[idx] = val
			search(idx + 1)
		}
	}

	search(0)
	return minPresses
}

func backSubstitute(aug [][]int, pivots []int, solution []int, nButtons int) bool {
	for i, col := range pivots {
		val := aug[i][nButtons]
		for j := 0; j < nButtons; j++ {
			if j != col {
				val -= aug[i][j] * solution[j]
			}
		}

		pivot := aug[i][col]
		if pivot == 0 {
			if val != 0 {
				return false
			}
			continue
		}
		if val%pivot != 0 || val/pivot < 0 {
			return false
		}
		solution[col] = val / pivot
	}
	return true
}

// GF(2) solver for Part 1
func solveGF2(matrix [][]int, target []int) int {
	if len(matrix) == 0 {
		return 0
	}

	aug := augment(matrix, target)
	pivots := rrefGF2(aug)

	// Check consistency
	for row := len(pivots); row < len(aug); row++ {
		if aug[row][len(matrix[0])] == 1 {
			return 0
		}
	}

	return searchMinGF2(aug, pivots, findFreeVars(pivots, len(matrix[0])))
}

func rrefGF2(aug [][]int) []int {
	nCols := len(aug[0]) - 1
	var pivots []int
	row := 0

	for col := 0; col < nCols; col++ {
		pivotRow := findPivot(aug, row, col, true)
		if pivotRow == -1 {
			continue
		}

		aug[row], aug[pivotRow] = aug[pivotRow], aug[row]
		pivots = append(pivots, col)

		for r := range aug {
			if r != row && aug[r][col] == 1 {
				for c := range aug[r] {
					aug[r][c] ^= aug[row][c]
				}
			}
		}
		row++
	}
	return pivots
}

func searchMinGF2(aug [][]int, pivots, freeVars []int) int {
	nCols := len(aug[0]) - 1
	minPresses := nCols + 1

	for mask := 0; mask < (1 << len(freeVars)); mask++ {
		solution := make([]int, nCols)
		for i, v := range freeVars {
			solution[v] = (mask >> i) & 1
		}

		// Back substitution in GF(2)
		for i := len(pivots) - 1; i >= 0; i-- {
			col := pivots[i]
			val := aug[i][nCols]
			for j := col + 1; j < nCols; j++ {
				val ^= aug[i][j] * solution[j]
			}
			solution[col] = val
		}

		minPresses = min(minPresses, sum(solution))
	}
	return minPresses
}

func findPivot(aug [][]int, startRow, col int, requireOne bool) int {
	for row := startRow; row < len(aug); row++ {
		if requireOne && aug[row][col] == 1 {
			return row
		}
		if !requireOne && aug[row][col] != 0 {
			return row
		}
	}
	return -1
}

func findFreeVars(pivots []int, nCols int) []int {
	pivotSet := make(map[int]bool, len(pivots))
	for _, p := range pivots {
		pivotSet[p] = true
	}

	var free []int
	for i := 0; i < nCols; i++ {
		if !pivotSet[i] {
			free = append(free, i)
		}
	}
	return free
}

func sum(arr []int) int {
	total := 0
	for _, v := range arr {
		total += v
	}
	return total
}
