/*
Day 8: Playground - Junction Box Circuit Problem

Part 1: Connect 1000 closest pairs, find product of 3 largest circuits
Part 2: Connect all boxes into one circuit, find product of last connection's X coordinates

Algorithm: Kruskal's MST with Union-Find
*/
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	x, y, z int
}

type Edge struct {
	i, j int
	dist float64
}

type UnionFind struct {
	parent []int
	size   []int
}

func NewUnionFind(n int) *UnionFind {
	uf := &UnionFind{
		parent: make([]int, n),
		size:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		uf.parent[i] = i
		uf.size[i] = 1
	}
	return uf
}

func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)

	if rootX == rootY {
		return false
	}

	if uf.size[rootX] < uf.size[rootY] {
		rootX, rootY = rootY, rootX
	}

	uf.parent[rootY] = rootX
	uf.size[rootX] += uf.size[rootY]
	return true
}

func (uf *UnionFind) GetCircuitSizes() []int {
	sizeMap := make(map[int]int)
	for i := 0; i < len(uf.parent); i++ {
		root := uf.Find(i)
		sizeMap[root] = uf.size[root]
	}

	sizes := make([]int, 0, len(sizeMap))
	for _, size := range sizeMap {
		sizes = append(sizes, size)
	}
	return sizes
}

func distance(p1, p2 Point) float64 {
	dx := float64(p1.x - p2.x)
	dy := float64(p1.y - p2.y)
	dz := float64(p1.z - p2.z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func parseInput(filename string) ([]Point, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	points := []Point{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) == 3 {
			x, _ := strconv.Atoi(parts[0])
			y, _ := strconv.Atoi(parts[1])
			z, _ := strconv.Atoi(parts[2])
			points = append(points, Point{x, y, z})
		}
	}
	return points, scanner.Err()
}

func buildEdges(points []Point) []Edge {
	n := len(points)
	edges := make([]Edge, 0, n*(n-1)/2)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dist := distance(points[i], points[j])
			edges = append(edges, Edge{i, j, dist})
		}
	}

	sort.Slice(edges, func(i, j int) bool {
		return edges[i].dist < edges[j].dist
	})

	return edges
}

func solvePart1(points []Point, edges []Edge) int {
	n := len(points)
	uf := NewUnionFind(n)

	processed := 0
	connected := 0

	for _, edge := range edges {
		if processed >= 1000 {
			break
		}
		processed++
		if uf.Union(edge.i, edge.j) {
			connected++
		}
	}

	sizes := uf.GetCircuitSizes()
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})

	fmt.Printf("Part 1: Processed %d pairs, made %d connections\n", processed, connected)
	fmt.Printf("Part 1: %d circuits remain\n", len(sizes))
	fmt.Printf("Part 1: Top circuit sizes: %v\n", sizes[:min(5, len(sizes))])

	if len(sizes) < 3 {
		return 0
	}
	return sizes[0] * sizes[1] * sizes[2]
}

func solvePart2(points []Point, edges []Edge) int {
	n := len(points)
	uf := NewUnionFind(n)
	numCircuits := n
	var lastI, lastJ int

	for _, edge := range edges {
		if uf.Union(edge.i, edge.j) {
			numCircuits--
			lastI, lastJ = edge.i, edge.j

			if numCircuits == 1 {
				break
			}
		}
	}

	fmt.Printf("Part 2: Last connection joins boxes %d and %d\n", lastI, lastJ)
	fmt.Printf("Part 2: Box %d at (%d, %d, %d)\n",
		lastI, points[lastI].x, points[lastI].y, points[lastI].z)
	fmt.Printf("Part 2: Box %d at (%d, %d, %d)\n",
		lastJ, points[lastJ].x, points[lastJ].y, points[lastJ].z)

	return points[lastI].x * points[lastJ].x
}

func main() {
	points, err := parseInput("input.txt")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	fmt.Printf("Loaded %d junction boxes\n\n", len(points))

	edges := buildEdges(points)
	fmt.Printf("Generated %d possible connections\n\n", len(edges))

	part1 := solvePart1(points, edges)
	fmt.Printf("Part 1 Answer: %d\n\n", part1)

	part2 := solvePart2(points, edges)
	fmt.Printf("Part 2 Answer: %d\n", part2)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
