package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("day1/input.txt")
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
	count := 0
	wheelPosition := 50

	fmt.Println(wheelPosition)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		instruction := scanner.Text()
		rotations := processInstruction(instruction, &wheelPosition)
		fmt.Println("Rotations:", rotations)
		fmt.Println("New wheel position:", wheelPosition)
		count += rotations
		if wheelPosition == 0 {
			count++
		}

		fmt.Println("Total rotations:", count)
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Total times at position 0:", count)
}

func processInstruction(instruction string, wheelPosition *int) int {
	// format L12 or R4
	direction := instruction[0]
	distance := instruction[1:]

	newPosition := *wheelPosition

	fmt.Println("current position:", newPosition)

	fullRotations := 0

	if newPosition == 0 && direction == 'L' {
		fullRotations -= 1
	}

	turns, rotations := parseDistance(distance)

	fullRotations += rotations

	switch direction {
	case 'L':
		newPosition -= turns
		fmt.Println("turn left", turns)
	case 'R':
		newPosition += turns
		fmt.Println("turn right", turns)
	}

	if newPosition < 0 {
		newPosition += 100
		fullRotations++
	} else if newPosition > 100 {
		newPosition -= 100
		fullRotations++
	} else if newPosition == 100 {
		newPosition = 0
	}

	*wheelPosition = newPosition

	return fullRotations
}

func parseDistance(distance string) (turns int, rotations int) {
	val, err := strconv.Atoi(distance)
	if err != nil {
		log.Fatal(err)
	}
	// must be between 0-99
	// remaining turns after full rotations
	turns = val % 100
	// number of full rotations
	rotations = val / 100

	return turns, rotations
}
