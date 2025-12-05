package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Range struct {
	start int
	end   int
}

func main() {
	file, err := os.Open("day5/input.txt")
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
	sum := 0
	isFresh := false
	ranges := []Range{}

	part1 := false

	for scanner.Scan() {
		line = scanner.Text()
		if isRange(line) {
			start, end, err := parseRange(line)
			if err != nil {
				log.Fatal(err)
			}
			// add to ranges
			newRange := Range{start: start, end: end}
			ranges = append(ranges, newRange)
			fmt.Printf("Adding fresh ingredients from %d to %d\n", start, end)
		} else if part1 {
			if line == "" {
				// skip blank lines
				continue
			} else {
				// single ingredient - check if fresh
				ingredient, err := strconv.Atoi(line)
				if err != nil {
					log.Fatal(err)
				}
				isFresh = isIngredientFresh(ingredient, ranges)
				if isFresh {
					fmt.Printf("Ingredient %d is fresh\n", ingredient)
					sum++
				} else {
					fmt.Printf("Ingredient %d is NOT fresh\n", ingredient)
				}
			}
		} else {
			break
		}

		ranges = mergeOverlappingRanges(ranges)

		fmt.Println("Parsed ranges:", ranges)
		if !part1 {
			sum = getTotalFreshIngredients(ranges)
		}
		fmt.Println("Total fresh ingredients so far:", sum)

		// ingredients are fresh if they fall into a range in the input
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Final parsed ranges:", ranges)
	fmt.Println("Total fresh ingredients:", sum)
}

func isRange(s string) bool {
	return strings.Contains(s, "-")
}

func parseRange(s string) (int, int, error) {
	var start, end int
	_, err := fmt.Sscanf(s, "%d-%d", &start, &end)
	if err != nil {
		return 0, 0, err
	}
	return start, end, nil
}

func isIngredientFresh(ingredient int, ranges []Range) bool {
	for _, r := range ranges {
		if ingredient >= r.start && ingredient <= r.end {
			return true
		}
	}
	return false
}

func mergeOverlappingRanges(ranges []Range) []Range {
	// sort by start
	if len(ranges) == 0 {
		return ranges
	}
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].start < ranges[j].start
	})
	merged := []Range{ranges[0]}
	for i := 1; i < len(ranges); i++ {
		last := &merged[len(merged)-1]
		current := ranges[i]
		if current.start <= last.end+1 {
			// overlap or contiguous - merge
			last.end = max(last.end, current.end)
		} else {
			// no overlap - add to merged
			merged = append(merged, current)
		}
	}
	return merged
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getTotalFreshIngredients(ranges []Range) int {
	total := 0
	for _, r := range ranges {
		total += r.end - r.start + 1
	}
	return total
}
