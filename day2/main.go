package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("day2/input.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	sum := 0

	for _, record := range records {
		fmt.Println(record)
		for i := 0; i < len(record); i++ {
			range1, range2 := parseRange(record[i])
			for id := range1; id <= range2; id++ {
				if isInvalidId(id) {
					fmt.Printf("Invalid ID: %d\n", id)
					sum += id
				}
			}
		}
	}
	fmt.Printf("Sum of invalid IDs: %d\n", sum)
}

func isInvalidId(id int) bool {
	// an ID is invalid if it is made only of some sequence of digits repeated at least twice
	idStr := strconv.Itoa(id)
	// a sequence repeated x number of times is invalid
	// LOOK AT SPLITTING AND REPEATING
	n := len(idStr)
	for seqLen := 1; seqLen <= n/2; seqLen++ {
		// if n is not divisible by seqLen, skip
		if n%seqLen != 0 {
			continue
		}
		seq := idStr[:seqLen]
		repeated := strings.Repeat(seq, n/seqLen)
		if repeated == idStr {
			return true // Invalid: made up of repeated sequence
		}
	}
	return false
}

func parseRange(rng string) (int, int) {
	//split on -
	parts := strings.Split(rng, "-")
	if len(parts) != 2 {
		log.Fatalf("invalid range: %s", rng)
	}
	start := parts[0]
	end := parts[1]
	startInt, err := strconv.Atoi(start)
	if err != nil {
		log.Fatalf("invalid range start: %s", start)
	}
	endInt, err := strconv.Atoi(end)
	if err != nil {
		log.Fatalf("invalid range end: %s", end)
	}
	return startInt, endInt
}
