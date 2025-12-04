package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const roll = "@"
const removed = "x"

func main() {
	file, err := os.Open("day4/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)

	line := ""

	grid := make([][]string, 0)
	sum := 0

	for scanner.Scan() {
		line = scanner.Text()
		row, err := parseLine(line)
		if err != nil {
			log.Fatal(err)
		}
		grid = append(grid, row)
		fmt.Println(row)
	}

	// remove as many rolls as possible
	// recurse until no more rolls can be removed
	removedThisRound := 0
	for {
		removedThisRound = removeRolls(grid)
		sum += removedThisRound
		if removedThisRound == 0 {
			break
		}
	}

	fmt.Println("Total rolls accessible by forklift:", sum)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func parseLine(s string) ([]string, error) {
	var symbols []string
	for _, r := range s {
		symbol := string(r)
		symbols = append(symbols, symbol)
	}
	return symbols, nil
}

func getAdjacentRolls(grid [][]string, row int, col int) [][2]int {
	var adjacentRolls [][2]int
	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, d := range directions {
		newRow := row + d[0]
		newCol := col + d[1]
		if newRow < 0 || newCol < 0 || newRow >= len(grid) || newCol >= len(grid[0]) {
			continue
		}
		if grid[newRow][newCol] != roll {
			continue
		}
		adjacentRolls = append(adjacentRolls, [2]int{newRow, newCol})
	}
	return adjacentRolls
}

func removeRoll(grid [][]string, row int, col int) {
	grid[row][col] = removed
}

func printGrid(grid [][]string) {
	for _, row := range grid {
		fmt.Println(row)
	}
}

func removeRolls(grid [][]string) int {
	sum := 0
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[0]); c++ {
			if grid[r][c] == roll {
				adjacentRolls := getAdjacentRolls(grid, r, c)
				fmt.Printf("Roll at (%d, %d) has adjacent rolls at: %v\n", r, c, adjacentRolls)
				if len(adjacentRolls) < 4 {
					removeRoll(grid, r, c)
					fmt.Println("Removed roll at:", r, c)
					fmt.Println("New grid:")
					printGrid(grid)
					sum++
				} else {
					fmt.Println("Forklift cannot access this roll.")
				}
			}
		}
	}
	return sum
}
