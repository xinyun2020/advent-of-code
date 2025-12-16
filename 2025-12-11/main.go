/*
Day 11: Reactor

Part 1: Count all paths from "you" to "out" in a directed graph
Part 2: Count paths from "svr" to "out" that visit both "dac" and "fft"
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Graph map[string][]string

type memoKey struct {
	node    string
	visited int
}

func main() {
	graph, err := parseInput("input.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Part 1: %d\n", graph.CountPaths("you", "out", nil))
	fmt.Printf("Part 2: %d\n", graph.CountPaths("svr", "out", []string{"dac", "fft"}))
}

func parseInput(filename string) (Graph, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	graph := make(Graph)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		from := strings.TrimSpace(parts[0])
		graph[from] = strings.Fields(parts[1])
	}

	return graph, scanner.Err()
}

// CountPaths counts paths from start to end.
// If required is non-empty, only counts paths visiting all required nodes.
func (g Graph) CountPaths(start, end string, required []string) int {
	reqIndex := make(map[string]int, len(required))
	for i, r := range required {
		reqIndex[r] = i
	}
	allVisited := (1 << len(required)) - 1

	memo := make(map[memoKey]int)

	var dfs func(node string, visited int) int
	dfs = func(node string, visited int) int {
		if idx, ok := reqIndex[node]; ok {
			visited |= 1 << idx
		}

		if node == end {
			if visited == allVisited {
				return 1
			}
			return 0
		}

		key := memoKey{node, visited}
		if val, ok := memo[key]; ok {
			return val
		}

		total := 0
		for _, next := range g[node] {
			total += dfs(next, visited)
		}

		memo[key] = total
		return total
	}

	return dfs(start, 0)
}
