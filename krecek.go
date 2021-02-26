package main

import "fmt"

const (
	WALL = '#'
)

func canGo(m string, pos int, h, w int) (can []int) {
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
	if a == -1 || b == -1 {
		return orig
	} else if orig == -1 {
		return a + b
	} else if orig > (a + b) {
		return a + b
	} else {
		return orig
	}
}

func FW(matrix [][]int) {
	for k := 0; k < len(matrix); k++ {
		for i := 0; i < len(matrix); i++ {
			for j := 0; j < len(matrix); j++ {
				matrix[i][j] = min(matrix[i][j], matrix[i][k], matrix[k][j])
			}
		}
	}
}

func makeDistances(maps []string, h, w int) (distances [][][]int) {
	distances = make([][][]int, len(maps))
	for i, m := range maps {
		distances[i] = make([][]int, len(m))
		for j := range distances[i] {
			distances[i][j] = make([]int, len(m))
			for k := range distances[i][j] { // init by -1
				distances[i][j][k] = -1
			}
		}
		for j := range distances[i] {
			can := canGo(m, j, h, w)
			for _, c := range can {
				distances[i][j][c] = 1
			}
		}
		for j := 0; j < len(distances[i]); j++ {
			distances[i][j][j] = 0
		}
		FW(distances[i])
	}
	return
}

func shortestWay(floor [][]int, start, end, last, stepsBefore int) int {
	if stepsBefore == -1 || floor[start][end] == -1 {
		return last
	} else if last == -1 {
		return floor[start][end] + stepsBefore
	} else if last < floor[start][end]+stepsBefore {
		return last
	} else {
		return floor[start][end] + stepsBefore
	}
}

func wayDown(maze []*[][]int, stairs [][]int, krecek, floor, lowestKrecek int) []int {
	var AFloor, FFloor []int

	if floor == lowestKrecek { // at lowest floor
		AFloor = make([]int, len(stairs[floor-1]))
		for i, e := range stairs[floor-1] {
			AFloor[i] = (*maze[floor])[krecek][e]
		}
		return AFloor
	}

	AFloor = make([]int, len(stairs[floor])) // goto stairs
	for i, e := range stairs[floor] {
		AFloor[i] = (*maze[floor])[krecek][e]
		if AFloor[i] != -1 {
			AFloor[i]++
		}
	}

	for floorNum := range stairs[floor : lowestKrecek-1] {
		fmt.Println(AFloor)
		f := floorNum + floor
		FFloor = make([]int, len(stairs[f+1]))
		for i := range FFloor { // init FFloor by -1
			FFloor[i] = -1
		}
		fmt.Println("maze:", f+1, stairs[f][0], (*maze[f+1])[stairs[f][0]])
		for i, start := range stairs[f] { // try goto stairs in lower floor
			for j, end := range stairs[f] {
				AFloor[j] = shortestWay(*maze[f+1], start, end, AFloor[j], AFloor[i])
			}
		}
		for i, start := range stairs[f] { // go to stairs down
			for j, end := range stairs[f+1] {
				FFloor[j] = shortestWay(*maze[f+1], start, end, FFloor[j], AFloor[i])
			}
		}
		for i := range FFloor {
			if FFloor[i] != -1 {
				FFloor[i]++
			}
		}
		AFloor = FFloor
		fmt.Println(AFloor)
	}
	return AFloor
}

func saveKrecky(maze []*[][]int, stairs [][]int, krecci [][]int, lowestKrecek, pocetKrecku int) int {
	var whenArrive [][]int
	var krecek int
	var minDistance int = -1
	whenArrive = make([][]int, pocetKrecku)
	for i := range whenArrive {
		whenArrive[i] = make([]int, len(stairs[lowestKrecek-1])+len(krecci[lowestKrecek]))
		for j := range whenArrive[i] {
			whenArrive[i][j] = -1
		}
	}

	for f, KFloor := range krecci {
		for _, k := range KFloor {
			arrive := wayDown(maze, stairs, k, f, lowestKrecek)
			//fmt.Println(arrive)
			if f != lowestKrecek { // normal floor
				for i, w := range arrive {
					if w != -1 {
						for j := range stairs[lowestKrecek-1] { // try go to stairs
							whenArrive[krecek][j] = shortestWay(*maze[lowestKrecek], stairs[lowestKrecek-1][i], stairs[lowestKrecek-1][j], min(arrive[j], whenArrive[krecek][j], 0), w)
						}
						for j := range krecci[lowestKrecek] { // goto krecek
							whenArrive[krecek][j+len(stairs[lowestKrecek-1])] = shortestWay(*maze[lowestKrecek], stairs[lowestKrecek-1][i], krecci[lowestKrecek][j], whenArrive[krecek][j+len(stairs[lowestKrecek-1])], w)
						}
					}
				}
			} else { // lowest floor
				for j, w := range arrive {
					whenArrive[krecek][j] = w
				}
				for j := range krecci[lowestKrecek] { // goto krecek
					whenArrive[krecek][j+len(stairs[lowestKrecek-1])] = shortestWay(*maze[lowestKrecek], k, krecci[lowestKrecek][j], whenArrive[krecek][j+len(stairs[lowestKrecek-1])], 0)
				}
			}
			krecek++
		}
	}
	//fmt.Println(whenArrive)
	for i := 0; i < len(whenArrive[0]); i++ { // find shortest way
		var sum int = -1
		for j := 0; j < len(whenArrive); j++ {
			if whenArrive[j][i] == -1 {
				sum = -1
				break
			} else {
				switch sum {
				case -1:
					sum = whenArrive[j][i]
				default:
					sum += whenArrive[j][i]
				}
			}
		}
		if sum != -1 && (sum < minDistance || minDistance == -1) {
			minDistance = sum
		}
	}
	return minDistance
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

	for i := 0; i < Q; i++ { // load questions
		var floors []*[][]int
		var floorsC int
		var staredFloors int
		var stairs [][]int
		var pocetKrecku int
		var krecci [][]int
		var mapTofloor map[int]int
		var lowestKrecek int

		fmt.Scan(&floorsC)
		floors = make([]*[][]int, floorsC)
		mapTofloor = make(map[int]int)
		for j := 0; j < floorsC; j++ { // load floors sequence
			var f int
			fmt.Scan(&f)
			floors[j] = &distances[f]
			mapTofloor[f] = j // puvodni cislo podlazi na cislo v dotazu
		}

		fmt.Scan(&staredFloors)
		stairs = make([][]int, staredFloors)
		for j := 0; j < staredFloors; j++ { // load stairs
			var startsCount int
			fmt.Scan(&startsCount)
			for k := 0; k < startsCount; k++ {
				var x, y int
				fmt.Scan(&y, &x)
				stairs[j] = append(stairs[j], y*w+x)
			}
		}

		fmt.Scan(&pocetKrecku)
		krecci = make([][]int, floorsC)
		for j := 0; j < pocetKrecku; j++ { // load krecky
			var floor, x, y int
			fmt.Scan(&floor, &y, &x)
			krecci[mapTofloor[floor]] = append(krecci[mapTofloor[floor]], y*w+x)
			if lowestKrecek < mapTofloor[floor] { // ktery krecek je nejnize
				lowestKrecek = mapTofloor[floor]
			}
		}
		distances := saveKrecky(floors, stairs, krecci, lowestKrecek, pocetKrecku)
		if distances == -1 {
			fmt.Println("social distancing")
		} else {
			fmt.Println(distances)
		}
	}
}
