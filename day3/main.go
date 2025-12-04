package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

//redo if have time - i dont like the way i handled slices and copying

func main() {
	file, err := os.Open("day3/input.txt")
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

	// get largest 2 values in string
	// check order
	// biggest number from left to right
	line := ""
	sum := 0

	for scanner.Scan() {
		line = scanner.Text()
		fmt.Println(line)
		digits, err := stringToDigitSlice(line)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(digits)

		largest := findLargestJoltage(digits, 12)
		sum += largest

		fmt.Printf("Found largest number: %d\n", largest)
		fmt.Println("--------------------")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sum of largest joltages:", sum)
}
func stringToDigitSlice(s string) ([]int, error) {
	var digits []int
	for _, r := range s {
		digit, err := strconv.Atoi(string(r))
		if err != nil {
			return nil, err
		}
		digits = append(digits, digit)
	}
	return digits, nil
}

func sortDigitsDescending(digits []int) []int {
	sort.Slice(digits, func(i, j int) bool {
		return digits[i] > digits[j]
	})
	return digits
}

func mapDigitsToIndices(digits []int) map[int][]int {
	digitPositions := make(map[int][]int)

	for index, digit := range digits {
		digitPositions[digit] = append(digitPositions[digit], index)
	}

	return digitPositions
}

func findLargestJoltage(digits []int, numDigits int) int {
	unsortedDigits := make([]int, len(digits))
	copy(unsortedDigits, digits)
	// first digit is largest within the first len(digits) - numDigits + 1
	firstValue, firstIndex := findLargestDigit(digits, 0, len(digits)-numDigits+1)
	fmt.Println("First digit:", firstValue, "at index", firstIndex)

	// sort digits
	mapped := mapDigitsToIndices(unsortedDigits)
	sortDigitsDescending(digits)
	fmt.Println("Sorted digits:", digits)
	fmt.Println("Mapped digits to indices:", mapped)

	// Initialize result as an empty slice
	result := make([]int, 0, numDigits)
	result = append(result, firstValue)

	// recurse for rest of list
	currentIndex := firstIndex + 1
	for i := 1; i < numDigits; i++ {
		fmt.Println("Current index:", currentIndex)
		fmt.Println("i:", i)
		endIndex := len(digits) - (numDigits - i - 1)
		nextValue, nextIndex := findLargestDigit(unsortedDigits, currentIndex, endIndex)
		result = append(result, nextValue)
		currentIndex = nextIndex + 1
	}

	// Convert result slice to a single number
	largestNumber := 0
	for _, digit := range result {
		largestNumber = largestNumber*10 + digit
	}
	return largestNumber
}

func findLargestDigit(digits []int, startIndex, endIndex int) (int, int) {
	copyDigits := make([]int, len(digits))
	copy(copyDigits, digits)
	// take slice from 0 to endIndex
	fmt.Println("startIndex:", startIndex)
	fmt.Println("endIndex:", endIndex)
	fmt.Println("digits:", digits)
	subSlice := copyDigits[startIndex:endIndex]
	fmt.Println("Subslice:", subSlice)
	mapped := mapDigitsToIndices(subSlice)
	sorted := sortDigitsDescending(subSlice)

	largest := sorted[0]
	index := mapped[largest][0] + startIndex
	fmt.Println("largestDigit in subslice:", largest)
	return largest, index
}
