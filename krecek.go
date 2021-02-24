package main

import "fmt"

const (
	WALL = '#'
)

func canGo(m string, pos int, h, w int) (can []int) {
	if m[pos] == WALL {
		return
	}
	can = make([]int, 0)
	withWalls := make([]int, 0)

	x := pos % w
	y := pos / w
	if x > 0 {
		withWalls = append(withWalls, y*w+x-1)
	}
	if y > 0 {
		withWalls = append(withWalls, (y-1)*w+x)
	}
	if x < (w - 1) {
		withWalls = append(withWalls, y*w+x+1)
	}
	if y < (h - 1) {
		withWalls = append(withWalls, (y+1)*w+x)
	}
	for _, p := range withWalls {
		if m[p] != WALL {
			can = append(can, p)
		}
	}
	return
}

func min(orig, a, b int) int {
	if a == 0 || b == 0 {
		return orig
	} else if orig == 0 {
		return a + b
	} else if orig > (a + b) {
		return a + b
	} else {
		return orig
	}
}

func FW(matrix [][]int) [][]int {
	for k := 0; k < len(matrix); k++ {
		for i := 0; i < len(matrix); i++ {
			for j := 0; j < len(matrix); j++ {
				matrix[i][j] = min(matrix[i][j], matrix[i][k], matrix[k][j])
			}
		}
	}
	for i := 0; i < len(matrix); i++ {
		matrix[i][i] = -1
	}
	return matrix
}

func makeDistances(maps []string, h, w int) (distances [][][]int) {
	distances = make([][][]int, len(maps))
	for i, m := range maps {
		distances[i] = make([][]int, len(m))
		for j, _ := range distances[i] {
			distances[i][j] = make([]int, len(m))
		}
		for j, _ := range distances[i] {
			can := canGo(m, j, h, w)
			for _, c := range can {
				distances[i][j][c] = 1
			}
		}
		distances[i] = FW(distances[i])
	}
	return
}

func main() {
	var Q, F, h, w int
	var maps []string
	var distances [][][]int

	fmt.Scan(&Q, &F, &h, &w)
	maps = make([]string, F)
	for i := 0; i < F; i++ { // load map
		for j := 0; j < h; j++ {
			var tmp string
			fmt.Scan(&tmp)
			maps[i] += tmp
		}
	}
	distances = makeDistances(maps, h, w) // compute distances matrix
	for i := 0; i < Q; i++ {              // load questions
		var floors []*[][]int
		var floorsC int
		var staredFloors int
		var stars [][]int

		fmt.Scan(&floorsC)
		floors = make([]*[][]int, floorsC)
		for j := 0; j < floorsC; j++ { // load floors sequence
			var f int
			fmt.Scan(&f)
			floors = append(floors, &distances[f])
		}

		fmt.Scan(&staredFloors)
		stars = make([][]int, staredFloors)
		for j := 0; j < staredFloors; j++ { // load stars
			var startsCount int
			fmt.Scan(&startsCount)
			for k := 0; k < startsCount; k++ {
				var x, y int
				fmt.Scan(&y, &x)
				stars[j] = append(stars[j], y*w+x)
			}
		}
	}
}
