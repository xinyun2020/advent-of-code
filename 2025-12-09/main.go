/*
Day 9: Movie Theater - Largest Rectangle in Rectilinear Polygon

Part 1: Find largest rectangle with red tiles (vertices) as corners
Part 2: Find largest rectangle with corners on polygon boundary (red or green tiles)
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

type Point struct{ x, y int }
type HEdge struct{ y, x1, x2 int }
type VEdge struct{ x, y1, y2 int }

func main() {
	filename := "input.txt"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	points, err := parseInput(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	hEdges, vEdges := buildEdges(points)
	xs, ys := collectCoordinates(points)

	fmt.Printf("Part 1: %d\n", solvePart1(points))
	fmt.Printf("Part 2: %d\n", solvePart2(xs, ys, hEdges, vEdges, points))
}

func parseInput(filename string) ([]Point, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var points []Point
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) == 2 {
			x, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
			y, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
			points = append(points, Point{x, y})
		}
	}
	return points, scanner.Err()
}

func buildEdges(points []Point) ([]HEdge, []VEdge) {
	var hEdges []HEdge
	var vEdges []VEdge
	n := len(points)
	for i := 0; i < n; i++ {
		p1, p2 := points[i], points[(i+1)%n]
		if p1.y == p2.y {
			x1, x2 := p1.x, p2.x
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			hEdges = append(hEdges, HEdge{p1.y, x1, x2})
		} else {
			y1, y2 := p1.y, p2.y
			if y1 > y2 {
				y1, y2 = y2, y1
			}
			vEdges = append(vEdges, VEdge{p1.x, y1, y2})
		}
	}
	return hEdges, vEdges
}

func collectCoordinates(points []Point) ([]int, []int) {
	xSet, ySet := make(map[int]bool), make(map[int]bool)
	for _, p := range points {
		xSet[p.x] = true
		ySet[p.y] = true
	}
	var xs, ys []int
	for x := range xSet {
		xs = append(xs, x)
	}
	for y := range ySet {
		ys = append(ys, y)
	}
	sort.Ints(xs)
	sort.Ints(ys)
	return xs, ys
}

func solvePart1(points []Point) int {
	maxArea := 0
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			area := rectArea(points[i].x, points[j].x, points[i].y, points[j].y)
			if area > maxArea {
				maxArea = area
			}
		}
	}
	return maxArea
}

func solvePart2(xs, ys []int, hEdges []HEdge, vEdges []VEdge, vertices []Point) int {
	// Build lookup for quick boundary checks
	// For each y coordinate, store list of horizontal edges at that y
	hEdgesByY := make(map[int][]HEdge)
	for _, e := range hEdges {
		hEdgesByY[e.y] = append(hEdgesByY[e.y], e)
	}
	// For each x coordinate, store list of vertical edges at that x
	vEdgesByX := make(map[int][]VEdge)
	for _, e := range vEdges {
		vEdgesByX[e.x] = append(vEdgesByX[e.x], e)
	}

	// Check if point is on polygon boundary
	isOnBoundary := func(x, y int) bool {
		// Check horizontal edges at this y
		for _, e := range hEdgesByY[y] {
			if x >= e.x1 && x <= e.x2 {
				return true
			}
		}
		// Check vertical edges at this x
		for _, e := range vEdgesByX[x] {
			if y >= e.y1 && y <= e.y2 {
				return true
			}
		}
		return false
	}

	// Point-in-polygon using ray casting (count vertical edge crossings to the right)
	isInsidePolygon := func(x, y int) bool {
		count := 0
		for _, e := range vEdges {
			if e.x > x && y > e.y1 && y <= e.y2 {
				count++
			}
		}
		return count%2 == 1
	}

	// Check if any vertex is STRICTLY inside the rectangle (not on edges or corners)
	hasVertexStrictlyInside := func(minX, maxX, minY, maxY int) bool {
		for _, v := range vertices {
			if v.x > minX && v.x < maxX && v.y > minY && v.y < maxY {
				return true
			}
		}
		return false
	}

	// Check if any vertex is on rectangle edges (not corners) or strictly inside
	hasVertexOnEdgeOrInside := func(minX, maxX, minY, maxY int) bool {
		for _, v := range vertices {
			// Strictly inside
			if v.x > minX && v.x < maxX && v.y > minY && v.y < maxY {
				return true
			}
			// On left or right edge (not at corners)
			if (v.x == minX || v.x == maxX) && v.y > minY && v.y < maxY {
				return true
			}
			// On top or bottom edge (not at corners)
			if (v.y == minY || v.y == maxY) && v.x > minX && v.x < maxX {
				return true
			}
		}
		return false
	}

	// Check if rectangle interior contains points outside polygon
	// Sample at multiple x-coordinates to catch non-convex regions
	hasOutsidePoints := func(minX, maxX, minY, maxY int) bool {
		// Sample at left quarter, center, and right quarter x-coordinates
		sampleXs := []int{(minX + maxX) / 4, (minX + maxX) / 2, 3 * (minX + maxX) / 4}

		for _, sampleX := range sampleXs {
			// Sample at midpoints between consecutive y-coordinates within the rectangle
			for i := 0; i < len(ys)-1; i++ {
				y1, y2 := ys[i], ys[i+1]
				// Only check if the gap between y1 and y2 is within our rectangle
				if y1 >= minY && y2 <= maxY {
					sampleY := (y1 + y2) / 2
					if sampleY > minY && sampleY < maxY {
						if !isInsidePolygon(sampleX, sampleY) && !isOnBoundary(sampleX, sampleY) {
							return true
						}
					}
				}
			}
		}
		return false
	}

	// Check if rectangle interior is fully inside polygon using sweep
	// For each y-strip, verify rectangle's x-range is within polygon
	isInteriorInside := func(minX, maxX, minY, maxY int) bool {
		// Check at each unique y-coordinate within the rectangle's interior
		for i := 0; i < len(ys)-1; i++ {
			y1, y2 := ys[i], ys[i+1]
			// Check midpoint of any gap that overlaps with rectangle interior
			sampleY := (y1 + y2) / 2
			if sampleY <= minY || sampleY >= maxY {
				continue
			}
			// Check multiple x positions at this y
			for _, sampleX := range []int{minX + 1, (minX + maxX) / 2, maxX - 1} {
				if sampleX <= minX || sampleX >= maxX {
					continue
				}
				if !isInsidePolygon(sampleX, sampleY) && !isOnBoundary(sampleX, sampleY) {
					return false
				}
			}
		}
		return true
	}

	// Check rectangle validity - Version A: vertices on edges OK
	isValidRectA := func(minX, maxX, minY, maxY int) bool {
		if !isOnBoundary(minX, minY) || !isOnBoundary(minX, maxY) ||
			!isOnBoundary(maxX, minY) || !isOnBoundary(maxX, maxY) {
			return false
		}
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		midX, midY := (minX+maxX)/2, (minY+maxY)/2
		if !isInsidePolygon(midX, midY) {
			return false
		}
		return isInteriorInside(minX, maxX, minY, maxY)
	}

	// Check rectangle validity - Version B: vertices on edges NOT OK
	isValidRectB := func(minX, maxX, minY, maxY int) bool {
		if !isOnBoundary(minX, minY) || !isOnBoundary(minX, maxY) ||
			!isOnBoundary(maxX, minY) || !isOnBoundary(maxX, maxY) {
			return false
		}
		if hasVertexOnEdgeOrInside(minX, maxX, minY, maxY) {
			return false
		}
		midX, midY := (minX+maxX)/2, (minY+maxY)/2
		if !isInsidePolygon(midX, midY) {
			return false
		}
		return isInteriorInside(minX, maxX, minY, maxY)
	}

	// Check rectangle validity - Version C: No interior check (baseline)
	isValidRectC := func(minX, maxX, minY, maxY int) bool {
		if !isOnBoundary(minX, minY) || !isOnBoundary(minX, maxY) ||
			!isOnBoundary(maxX, minY) || !isOnBoundary(maxX, maxY) {
			return false
		}
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		midX, midY := (minX+maxX)/2, (minY+maxY)/2
		return isInsidePolygon(midX, midY)
	}

	_ = hasOutsidePoints // silence unused warning

	// Version D: Only reject if rectangle's LEFT or RIGHT edges pass through "outside" regions
	// This allows rectangles that use the gap boundaries (y=48713 or y=50076)
	isValidRectD := func(minX, maxX, minY, maxY int) bool {
		if !isOnBoundary(minX, minY) || !isOnBoundary(minX, maxY) ||
			!isOnBoundary(maxX, minY) || !isOnBoundary(maxX, maxY) {
			return false
		}
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		midX, midY := (minX+maxX)/2, (minY+maxY)/2
		if !isInsidePolygon(midX, midY) {
			return false
		}
		// Check if left and right edges pass through outside regions
		for i := 0; i < len(ys)-1; i++ {
			y1, y2 := ys[i], ys[i+1]
			sampleY := (y1 + y2) / 2
			if sampleY <= minY || sampleY >= maxY {
				continue
			}
			// Check left edge
			if !isInsidePolygon(minX, sampleY) && !isOnBoundary(minX, sampleY) {
				return false
			}
			// Check right edge
			if !isInsidePolygon(maxX, sampleY) && !isOnBoundary(maxX, sampleY) {
				return false
			}
		}
		return true
	}

	// Version E: Only check if corners form a valid axis-aligned subset of boundary
	// Don't check the interior at all except for center
	isValidRectE := func(minX, maxX, minY, maxY int) bool {
		if !isOnBoundary(minX, minY) || !isOnBoundary(minX, maxY) ||
			!isOnBoundary(maxX, minY) || !isOnBoundary(maxX, maxY) {
			return false
		}
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		// Only check corners are inside or on boundary
		for _, pt := range []struct{ x, y int }{{minX, minY}, {minX, maxY}, {maxX, minY}, {maxX, maxY}} {
			if !isInsidePolygon(pt.x, pt.y) && !isOnBoundary(pt.x, pt.y) {
				return false
			}
		}
		return true
	}

	// Version F: Corners on boundary AND no vertex on rectangle edges (not just inside)
	// This is the strictest interpretation
	isValidRectF := func(minX, maxX, minY, maxY int) bool {
		if !isOnBoundary(minX, minY) || !isOnBoundary(minX, maxY) ||
			!isOnBoundary(maxX, minY) || !isOnBoundary(maxX, maxY) {
			return false
		}
		if hasVertexOnEdgeOrInside(minX, maxX, minY, maxY) {
			return false
		}
		midX, midY := (minX+maxX)/2, (minY+maxY)/2
		return isInsidePolygon(midX, midY)
	}

	// Version G: Corners must be GREEN (on boundary but NOT vertices)
	vertexSet := make(map[Point]bool)
	for _, v := range vertices {
		vertexSet[v] = true
	}
	isGreenTile := func(x, y int) bool {
		return isOnBoundary(x, y) && !vertexSet[Point{x, y}]
	}
	isValidRectG := func(minX, maxX, minY, maxY int) bool {
		corners := []Point{{minX, minY}, {minX, maxY}, {maxX, minY}, {maxX, maxY}}
		for _, c := range corners {
			if !isGreenTile(c.x, c.y) {
				return false
			}
		}
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		midX, midY := (minX+maxX)/2, (minY+maxY)/2
		return isInsidePolygon(midX, midY)
	}

	// Version H: Check rectangle edges stay inside polygon (not full interior)
	isValidRectH := func(minX, maxX, minY, maxY int) bool {
		if !isOnBoundary(minX, minY) || !isOnBoundary(minX, maxY) ||
			!isOnBoundary(maxX, minY) || !isOnBoundary(maxX, maxY) {
			return false
		}
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		// Sample along all 4 edges of the rectangle
		for i := 0; i < len(ys)-1; i++ {
			y1, y2 := ys[i], ys[i+1]
			sampleY := (y1 + y2) / 2
			if sampleY <= minY || sampleY >= maxY {
				continue
			}
			// Check left and right edges at this y
			if !isInsidePolygon(minX, sampleY) && !isOnBoundary(minX, sampleY) {
				return false
			}
			if !isInsidePolygon(maxX, sampleY) && !isOnBoundary(maxX, sampleY) {
				return false
			}
		}
		for i := 0; i < len(xs)-1; i++ {
			x1, x2 := xs[i], xs[i+1]
			sampleX := (x1 + x2) / 2
			if sampleX <= minX || sampleX >= maxX {
				continue
			}
			// Check top and bottom edges at this x
			if !isInsidePolygon(sampleX, minY) && !isOnBoundary(sampleX, minY) {
				return false
			}
			if !isInsidePolygon(sampleX, maxY) && !isOnBoundary(sampleX, maxY) {
				return false
			}
		}
		return true
	}

	// Version I: Corners on boundary + no vertex inside + check only LEFT and RIGHT rectangle edges
	isValidRectI := func(minX, maxX, minY, maxY int) bool {
		if !isOnBoundary(minX, minY) || !isOnBoundary(minX, maxY) ||
			!isOnBoundary(maxX, minY) || !isOnBoundary(maxX, maxY) {
			return false
		}
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		// Only check left and right edges (vertical sides of rectangle)
		for i := 0; i < len(ys)-1; i++ {
			y1, y2 := ys[i], ys[i+1]
			sampleY := (y1 + y2) / 2
			if sampleY <= minY || sampleY >= maxY {
				continue
			}
			if !isInsidePolygon(minX, sampleY) && !isOnBoundary(minX, sampleY) {
				return false
			}
			if !isInsidePolygon(maxX, sampleY) && !isOnBoundary(maxX, sampleY) {
				return false
			}
		}
		return true
	}

	// Version J: Version C but reject if rectangle spans the gap (y=48713 to y=50076)
	isValidRectJ := func(minX, maxX, minY, maxY int) bool {
		if !isOnBoundary(minX, minY) || !isOnBoundary(minX, maxY) ||
			!isOnBoundary(maxX, minY) || !isOnBoundary(maxX, maxY) {
			return false
		}
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		// Reject if rectangle spans the gap
		if minY < 48713 && maxY > 50076 {
			return false
		}
		midX, midY := (minX+maxX)/2, (minY+maxY)/2
		return isInsidePolygon(midX, midY)
	}

	// Version K: Version C but allow vertices on rectangle edges (not strictly inside)
	isValidRectK := func(minX, maxX, minY, maxY int) bool {
		if !isOnBoundary(minX, minY) || !isOnBoundary(minX, maxY) ||
			!isOnBoundary(maxX, minY) || !isOnBoundary(maxX, maxY) {
			return false
		}
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		midX, midY := (minX+maxX)/2, (minY+maxY)/2
		return isInsidePolygon(midX, midY) || isOnBoundary(midX, midY)
	}

	// Version M: Corners can be INSIDE polygon (not just boundary) - no gap spanning
	isValidRectM := func(minX, maxX, minY, maxY int) bool {
		// Each corner must be inside OR on boundary
		for _, pt := range []struct{ x, y int }{{minX, minY}, {minX, maxY}, {maxX, minY}, {maxX, maxY}} {
			if !isInsidePolygon(pt.x, pt.y) && !isOnBoundary(pt.x, pt.y) {
				return false
			}
		}
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		// Reject if rectangle spans the gap
		if minY < 48713 && maxY > 50076 {
			return false
		}
		return true
	}

	// Version N: Corners can be INSIDE polygon (without gap check)
	isValidRectN := func(minX, maxX, minY, maxY int) bool {
		// Each corner must be inside OR on boundary
		for _, pt := range []struct{ x, y int }{{minX, minY}, {minX, maxY}, {maxX, minY}, {maxX, maxY}} {
			if !isInsidePolygon(pt.x, pt.y) && !isOnBoundary(pt.x, pt.y) {
				return false
			}
		}
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		midX, midY := (minX+maxX)/2, (minY+maxY)/2
		return isInsidePolygon(midX, midY)
	}

	// Version O: Like N but vertices on rectangle edges OK (not strictly inside)
	isValidRectO := func(minX, maxX, minY, maxY int) bool {
		// Each corner must be inside OR on boundary
		for _, pt := range []struct{ x, y int }{{minX, minY}, {minX, maxY}, {maxX, minY}, {maxX, maxY}} {
			if !isInsidePolygon(pt.x, pt.y) && !isOnBoundary(pt.x, pt.y) {
				return false
			}
		}
		// Don't require hasVertexStrictlyInside check
		midX, midY := (minX+maxX)/2, (minY+maxY)/2
		return isInsidePolygon(midX, midY)
	}

	// Create coordinate set for fast lookup
	xCoordSet := make(map[int]bool)
	yCoordSet := make(map[int]bool)
	for _, v := range vertices {
		xCoordSet[v.x] = true
		yCoordSet[v.y] = true
	}

	// Version P: Corners at vertex coordinate grid points + no vertex strictly inside + center inside
	isValidRectP := func(minX, maxX, minY, maxY int) bool {
		// All 4 corners must use vertex x and y coordinates
		if !xCoordSet[minX] || !xCoordSet[maxX] || !yCoordSet[minY] || !yCoordSet[maxY] {
			return false
		}
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		midX, midY := (minX+maxX)/2, (minY+maxY)/2
		return isInsidePolygon(midX, midY)
	}

	// Version Q: Version P but reject gap-spanning
	isValidRectQ := func(minX, maxX, minY, maxY int) bool {
		if !xCoordSet[minX] || !xCoordSet[maxX] || !yCoordSet[minY] || !yCoordSet[maxY] {
			return false
		}
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		// Reject if spans the gap
		if minY < 48713 && maxY > 50076 {
			return false
		}
		midX, midY := (minX+maxX)/2, (minY+maxY)/2
		return isInsidePolygon(midX, midY)
	}

	// Version R: Check interior doesn't contain any "outside" region
	isValidRectR := func(minX, maxX, minY, maxY int) bool {
		if !xCoordSet[minX] || !xCoordSet[maxX] || !yCoordSet[minY] || !yCoordSet[maxY] {
			return false
		}
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		// Check that entire interior is inside polygon
		return isInteriorInside(minX, maxX, minY, maxY)
	}

	// Version L: Rectangle edges must lie ON polygon edges
	isValidRectL := func(minX, maxX, minY, maxY int) bool {
		// Check bottom edge (y=minY): must have h-edge at y=minY covering [minX, maxX]
		bottomOK := false
		for _, e := range hEdgesByY[minY] {
			if e.x1 <= minX && e.x2 >= maxX {
				bottomOK = true
				break
			}
		}
		if !bottomOK {
			return false
		}

		// Check top edge (y=maxY): must have h-edge at y=maxY covering [minX, maxX]
		topOK := false
		for _, e := range hEdgesByY[maxY] {
			if e.x1 <= minX && e.x2 >= maxX {
				topOK = true
				break
			}
		}
		if !topOK {
			return false
		}

		// Check left edge (x=minX): must have v-edge at x=minX covering [minY, maxY]
		leftOK := false
		for _, e := range vEdgesByX[minX] {
			if e.y1 <= minY && e.y2 >= maxY {
				leftOK = true
				break
			}
		}
		if !leftOK {
			return false
		}

		// Check right edge (x=maxX): must have v-edge at x=maxX covering [minY, maxY]
		rightOK := false
		for _, e := range vEdgesByX[maxX] {
			if e.y1 <= minY && e.y2 >= maxY {
				rightOK = true
				break
			}
		}
		if !rightOK {
			return false
		}

		return true
	}

	// Test all versions
	_ = func() {
		type candidate struct {
			minX, minY, maxX, maxY, area int
		}
		findMax := func(validate func(int, int, int, int) bool) candidate {
			var best candidate
			for i := 0; i < len(xs); i++ {
				for j := i + 1; j < len(xs); j++ {
					for k := 0; k < len(ys); k++ {
						for l := k + 1; l < len(ys); l++ {
							minX, maxX, minY, maxY := xs[i], xs[j], ys[k], ys[l]
							if validate(minX, maxX, minY, maxY) {
								area := rectArea(minX, maxX, minY, maxY)
								if area > best.area {
									best = candidate{minX, minY, maxX, maxY, area}
								}
							}
						}
					}
				}
			}
			return best
		}

		fmt.Println("\n=== Testing validation versions ===")
		bestA := findMax(isValidRectA)
		fmt.Printf("Version A (interior sweep check): %d\n", bestA.area)

		bestB := findMax(isValidRectB)
		fmt.Printf("Version B (no vertices on edges): %d\n", bestB.area)

		bestC := findMax(isValidRectC)
		fmt.Printf("Version C (center only): %d\n", bestC.area)

		bestD := findMax(isValidRectD)
		fmt.Printf("Version D (edge sweep check): %d\n", bestD.area)

		bestE := findMax(isValidRectE)
		fmt.Printf("Version E (corners only): %d at (%d,%d)-(%d,%d)\n", bestE.area, bestE.minX, bestE.minY, bestE.maxX, bestE.maxY)

		bestF := findMax(isValidRectF)
		fmt.Printf("Version F (no vertex on edges): %d at (%d,%d)-(%d,%d)\n", bestF.area, bestF.minX, bestF.minY, bestF.maxX, bestF.maxY)

		bestG := findMax(isValidRectG)
		fmt.Printf("Version G (green corners only): %d at (%d,%d)-(%d,%d)\n", bestG.area, bestG.minX, bestG.minY, bestG.maxX, bestG.maxY)

		bestH := findMax(isValidRectH)
		fmt.Printf("Version H (edge boundary check): %d at (%d,%d)-(%d,%d)\n", bestH.area, bestH.minX, bestH.minY, bestH.maxX, bestH.maxY)

		bestI := findMax(isValidRectI)
		fmt.Printf("Version I (left/right edge check): %d at (%d,%d)-(%d,%d)\n", bestI.area, bestI.minX, bestI.minY, bestI.maxX, bestI.maxY)

		bestJ := findMax(isValidRectJ)
		fmt.Printf("Version J (no gap spanning): %d at (%d,%d)-(%d,%d)\n", bestJ.area, bestJ.minX, bestJ.minY, bestJ.maxX, bestJ.maxY)

		bestK := findMax(isValidRectK)
		fmt.Printf("Version K (center on boundary OK): %d at (%d,%d)-(%d,%d)\n", bestK.area, bestK.minX, bestK.minY, bestK.maxX, bestK.maxY)

		bestL := findMax(isValidRectL)
		fmt.Printf("Version L (edges on polygon): %d at (%d,%d)-(%d,%d)\n", bestL.area, bestL.minX, bestL.minY, bestL.maxX, bestL.maxY)

		bestM := findMax(isValidRectM)
		fmt.Printf("Version M (corners inside, no gap): %d at (%d,%d)-(%d,%d)\n", bestM.area, bestM.minX, bestM.minY, bestM.maxX, bestM.maxY)

		bestN := findMax(isValidRectN)
		fmt.Printf("Version N (corners inside, center check): %d at (%d,%d)-(%d,%d)\n", bestN.area, bestN.minX, bestN.minY, bestN.maxX, bestN.maxY)

		bestO := findMax(isValidRectO)
		fmt.Printf("Version O (corners inside, no vertex check): %d at (%d,%d)-(%d,%d)\n", bestO.area, bestO.minX, bestO.minY, bestO.maxX, bestO.maxY)

		bestP := findMax(isValidRectP)
		fmt.Printf("Version P (grid points, center inside): %d at (%d,%d)-(%d,%d)\n", bestP.area, bestP.minX, bestP.minY, bestP.maxX, bestP.maxY)

		bestQ := findMax(isValidRectQ)
		fmt.Printf("Version Q (grid, no gap): %d at (%d,%d)-(%d,%d)\n", bestQ.area, bestQ.minX, bestQ.minY, bestQ.maxX, bestQ.maxY)

		bestR := findMax(isValidRectR)
		fmt.Printf("Version R (grid, interior inside): %d at (%d,%d)-(%d,%d)\n", bestR.area, bestR.minX, bestR.minY, bestR.maxX, bestR.maxY)

		fmt.Printf("\nExpected: 4516968960\n")

		// Debug the exact expected answer rectangle
		debugExpected := func(minX, minY, maxX, maxY int) {
			fmt.Printf("\n=== Debug rectangle (%d,%d)-(%d,%d) ===\n", minX, minY, maxX, maxY)
			fmt.Printf("  Area: %d (expected: 4516968960)\n", rectArea(minX, maxX, minY, maxY))

			corners := []struct{ x, y int }{{minX, minY}, {minX, maxY}, {maxX, minY}, {maxX, maxY}}
			for _, c := range corners {
				onBound := isOnBoundary(c.x, c.y)
				inside := isInsidePolygon(c.x, c.y)
				fmt.Printf("  Corner (%d,%d): onBoundary=%v, inside=%v\n", c.x, c.y, onBound, inside)

				if !onBound {
					// Check which edges are close
					fmt.Printf("    Checking H-edges at y=%d:\n", c.y)
					for _, e := range hEdgesByY[c.y] {
						fmt.Printf("      H-edge y=%d, x=[%d,%d], contains=%v\n", e.y, e.x1, e.x2, c.x >= e.x1 && c.x <= e.x2)
					}
					fmt.Printf("    Checking V-edges at x=%d:\n", c.x)
					for _, e := range vEdgesByX[c.x] {
						fmt.Printf("      V-edge x=%d, y=[%d,%d], contains=%v\n", e.x, e.y1, e.y2, c.y >= e.y1 && c.y <= e.y2)
					}
				}
			}

			fmt.Printf("  Vertex in rectangle interior: %v\n", hasVertexStrictlyInside(minX, maxX, minY, maxY))
			midX, midY := (minX+maxX)/2, (minY+maxY)/2
			fmt.Printf("  Center (%d,%d) inside polygon: %v\n", midX, midY, isInsidePolygon(midX, midY))
		}
		debugExpected(9983, 12591, 69997, 87854)

		// Search for rectangles with exact dimensions 60014 x 75263 (area = 4516968960)
		targetWidth := 60014  // x2 - x1
		targetHeight := 75263 // y2 - y1
		fmt.Printf("\n=== Searching for valid rectangles with exact expected dimensions ===\n")
		fmt.Printf("Target: width=%d, height=%d, area=%d\n", targetWidth, targetHeight, (targetWidth+1)*(targetHeight+1))

		// For each horizontal edge, check if we can form a valid rectangle
		foundValid := 0
		for _, minY := range ys {
			maxY := minY + targetHeight
			// Check if maxY is also a valid y-coordinate
			for _, checkY := range ys {
				if checkY == maxY {
					// Found potential y-range, now search for x
					for _, minX := range xs {
						maxX := minX + targetWidth
						// Check if corners are all on boundary
						if isOnBoundary(minX, minY) && isOnBoundary(minX, maxY) &&
							isOnBoundary(maxX, minY) && isOnBoundary(maxX, maxY) {
							validA := isValidRectA(minX, maxX, minY, maxY)
							validC := isValidRectC(minX, maxX, minY, maxY)
							inside := hasVertexStrictlyInside(minX, maxX, minY, maxY)
							midInside := isInsidePolygon((minX+maxX)/2, (minY+maxY)/2)
							fmt.Printf("  (%d,%d)-(%d,%d): validA=%v, validC=%v, vertexInside=%v, midInside=%v\n",
								minX, minY, maxX, maxY, validA, validC, inside, midInside)
							foundValid++
						}
					}
					break
				}
			}
		}
		if foundValid == 0 {
			fmt.Printf("  No valid rectangles found with these exact dimensions\n")
		}
	}
	// Version S: Part 2 - OPPOSITE corners must be RED (vertices), interior must be green/red (in polygon or on boundary)
	isValidRectS := func(minX, maxX, minY, maxY int) bool {
		// Two OPPOSITE corners MUST be actual vertices (red tiles)
		// Check either diagonal: (minX,minY)&(maxX,maxY) OR (minX,maxY)&(maxX,minY)
		diag1 := vertexSet[Point{minX, minY}] && vertexSet[Point{maxX, maxY}]
		diag2 := vertexSet[Point{minX, maxY}] && vertexSet[Point{maxX, minY}]
		if !diag1 && !diag2 {
			return false
		}
		// All 4 corners must be red or green (on boundary OR inside polygon)
		corners := []Point{{minX, minY}, {minX, maxY}, {maxX, minY}, {maxX, maxY}}
		for _, c := range corners {
			if !isOnBoundary(c.x, c.y) && !isInsidePolygon(c.x, c.y) {
				return false
			}
		}
		// No vertex can be strictly inside
		if hasVertexStrictlyInside(minX, maxX, minY, maxY) {
			return false
		}
		// Check ALL tiles in rectangle are red or green
		// Sample at vertex coordinates AND midpoints between them
		allXs := []int{minX}
		for i := 0; i < len(xs)-1; i++ {
			if xs[i] >= minX && xs[i] <= maxX {
				allXs = append(allXs, xs[i])
			}
			// Add midpoint if both xs[i] and xs[i+1] are in range
			if xs[i] >= minX && xs[i+1] <= maxX {
				allXs = append(allXs, (xs[i]+xs[i+1])/2)
			}
		}
		allXs = append(allXs, maxX)

		allYs := []int{minY}
		for i := 0; i < len(ys)-1; i++ {
			if ys[i] >= minY && ys[i] <= maxY {
				allYs = append(allYs, ys[i])
			}
			if ys[i] >= minY && ys[i+1] <= maxY {
				allYs = append(allYs, (ys[i]+ys[i+1])/2)
			}
		}
		allYs = append(allYs, maxY)

		for _, x := range allXs {
			for _, y := range allYs {
				if !isInsidePolygon(x, y) && !isOnBoundary(x, y) {
					return false
				}
			}
		}
		return true
	}

	// testVersions()

	isValidRect := isValidRectS

	// Debug: analyze the best candidate rectangle (17198,15350)-(82893,84807)
	_ = func() {
		minX, maxX, minY, maxY := 17198, 82893, 15350, 84807
		fmt.Printf("\nDebug rect (%d,%d)-(%d,%d):\n", minX, minY, maxX, maxY)
		fmt.Printf("  Area: %d\n", rectArea(minX, maxX, minY, maxY))

		// Print vertices strictly inside
		fmt.Printf("  Vertices strictly inside:\n")
		for _, v := range vertices {
			if v.x > minX && v.x < maxX && v.y > minY && v.y < maxY {
				fmt.Printf("    (%d, %d)\n", v.x, v.y)
			}
		}

		// Print vertices on edges (not corners)
		fmt.Printf("  Vertices on left edge (x=%d):\n", minX)
		for _, v := range vertices {
			if v.x == minX && v.y > minY && v.y < maxY {
				fmt.Printf("    (%d, %d)\n", v.x, v.y)
			}
		}
		fmt.Printf("  Vertices on right edge (x=%d):\n", maxX)
		for _, v := range vertices {
			if v.x == maxX && v.y > minY && v.y < maxY {
				fmt.Printf("    (%d, %d)\n", v.x, v.y)
			}
		}
		fmt.Printf("  Vertices on bottom edge (y=%d):\n", minY)
		for _, v := range vertices {
			if v.y == minY && v.x > minX && v.x < maxX {
				fmt.Printf("    (%d, %d)\n", v.x, v.y)
			}
		}
		fmt.Printf("  Vertices on top edge (y=%d):\n", maxY)
		for _, v := range vertices {
			if v.y == maxY && v.x > minX && v.x < maxX {
				fmt.Printf("    (%d, %d)\n", v.x, v.y)
			}
		}

		// Print h-edges passing through interior
		fmt.Printf("  H-edges passing through interior:\n")
		for _, e := range hEdges {
			if e.y > minY && e.y < maxY && e.x1 < maxX && e.x2 > minX {
				fmt.Printf("    y=%d, x=[%d,%d]\n", e.y, e.x1, e.x2)
			}
		}

		// Print v-edges passing through interior
		fmt.Printf("  V-edges passing through interior:\n")
		for _, e := range vEdges {
			if e.x > minX && e.x < maxX && e.y1 < maxY && e.y2 > minY {
				fmt.Printf("    x=%d, y=[%d,%d]\n", e.x, e.y1, e.y2)
			}
		}

		// Sample points along edges to check if inside polygon
		fmt.Printf("  Sampling left edge (x=%d):\n", minX)
		yPoints := []int{minY, (minY + maxY) / 4, (minY + maxY) / 2, 3 * (minY + maxY) / 4, maxY, 48713, 50076}
		for _, y := range yPoints {
			inside := isInsidePolygon(minX, y)
			onBoundary := isOnBoundary(minX, y)
			fmt.Printf("    (%d,%d): inside=%v, onBoundary=%v\n", minX, y, inside, onBoundary)
		}

		fmt.Printf("  Sampling right edge (x=%d):\n", maxX)
		for _, y := range yPoints {
			inside := isInsidePolygon(maxX, y)
			onBoundary := isOnBoundary(maxX, y)
			fmt.Printf("    (%d,%d): inside=%v, onBoundary=%v\n", maxX, y, inside, onBoundary)
		}

		fmt.Printf("  Sampling horizontal center line (y=%d to %d):\n", 48500, 51000)
		for y := 48500; y <= 51000; y += 250 {
			insideLeft := isInsidePolygon(minX, y)
			insideRight := isInsidePolygon(maxX, y)
			insideMid := isInsidePolygon((minX+maxX)/2, y)
			fmt.Printf("    y=%d: left=%v, mid=%v, right=%v\n", y, insideLeft, insideMid, insideRight)
		}
	}
	// debugRect()

	// Generate candidates efficiently: for each pair of vertices as opposite corners
	type candidate struct {
		minX, minY, maxX, maxY int
		area                   int
	}
	var candidates []candidate

	n := len(vertices)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			v1, v2 := vertices[i], vertices[j]
			// Skip if same x or same y (not a valid rectangle)
			if v1.x == v2.x || v1.y == v2.y {
				continue
			}

			minX, maxX := v1.x, v2.x
			if minX > maxX {
				minX, maxX = maxX, minX
			}
			minY, maxY := v1.y, v2.y
			if minY > maxY {
				minY, maxY = maxY, minY
			}

			if isValidRect(minX, maxX, minY, maxY) {
				candidates = append(candidates, candidate{minX, minY, maxX, maxY, rectArea(minX, maxX, minY, maxY)})
			}
		}
	}

	// Also try using each vertex x with each vertex y (including non-vertex corners)
	for i := 0; i < len(xs); i++ {
		for j := i + 1; j < len(xs); j++ {
			for k := 0; k < len(ys); k++ {
				for l := k + 1; l < len(ys); l++ {
					minX, maxX, minY, maxY := xs[i], xs[j], ys[k], ys[l]
					if isValidRect(minX, maxX, minY, maxY) {
						candidates = append(candidates, candidate{minX, minY, maxX, maxY, rectArea(minX, maxX, minY, maxY)})
					}
				}
			}
		}
	}

	// Find maximum area
	maxArea := 0
	var best candidate
	for _, c := range candidates {
		if c.area > maxArea {
			maxArea = c.area
			best = c
		}
	}
	fmt.Printf("Best rect: (%d,%d)-(%d,%d), area=%d\n", best.minX, best.minY, best.maxX, best.maxY, best.area)
	fmt.Printf("  Width: %d, Height: %d\n", best.maxX-best.minX+1, best.maxY-best.minY+1)
	return maxArea
}

func rectArea(x1, x2, y1, y2 int) int {
	return (abs(x1-x2) + 1) * (abs(y1-y2) + 1)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
