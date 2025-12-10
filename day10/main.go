package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize/convex/lp"
)

type MachineConfig struct {
	indicatorLights []bool
	buttonWirings   [][]bool
	joltages        []int
}

func main() {
	part1 := false
	buttonPresses := 0

	file, err := os.Open("day10/test.txt")
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
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
		machine := parseLine(line)
		// what is the fewest number of button presses to correctly configure all indicator lights to match?
		// all indicator lights start off
		best := -1
		if part1 {
			best = findBestButtonPresses(machine)
		} else {
			// Each machine needs to be configured to exactly the specified joltage levels to function properly
			// counters are all initially set to zero.
			// When you push a button, each listed counter is increased by 1.
			// fewest total presses required to correctly configure each machine's joltage level counters to match the specified joltage requirements.
			// only check joltages now
			best = findJoltageMatches(machine)
		}

		if best != -1 {
			buttonPresses += best
		}

		log.Println("-----")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Total button presses needed: %d\n", buttonPresses)
}

func parseLine(line string) MachineConfig {
	// indicator lights are in []
	// button wirings are in ()
	// joltages in {}
	// example line: [.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}
	indicatorLights := []bool{}
	buttonWirings := [][]int{}
	joltages := []int{}
	parts := strings.Fields(line)
	for _, part := range parts {
		log.Printf("Parsing line part: %v\n", part)
		if strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
			log.Println("Parsing indicator lights")
			// indicator lights
			indicators := part[1 : len(part)-1]
			for _, ch := range indicators {
				if ch == '#' {
					indicatorLights = append(indicatorLights, true)
				} else {
					indicatorLights = append(indicatorLights, false)
				}
			}
		} else if strings.HasPrefix(part, "(") && strings.HasSuffix(part, ")") {
			log.Println("Parsing button wirings")
			// button wirings
			wiringStr := part[1 : len(part)-1]
			wiringParts := strings.Split(wiringStr, ",")
			wiring := []int{}
			for _, wp := range wiringParts {
				num, err := strconv.Atoi(wp)
				if err != nil {
					log.Fatal(err)
				}
				wiring = append(wiring, num)
			}
			buttonWirings = append(buttonWirings, wiring)
		} else if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			log.Println("Parsing joltages")
			// joltages
			joltageStr := part[1 : len(part)-1]
			joltageParts := strings.Split(joltageStr, ",")
			for _, jp := range joltageParts {
				num, err := strconv.Atoi(jp)
				if err != nil {
					log.Fatal(err)
				}
				joltages = append(joltages, num)
			}
		}
	}
	return MachineConfig{
		indicatorLights: indicatorLights,
		buttonWirings:   parseWirings(buttonWirings, len(indicatorLights)),
		joltages:        joltages,
	}
}

func pressButton(buttonIndex int, indicatorLights []bool, buttonWirings [][]bool) {
	// apply button wiring transformation to indicator lights
	if buttonIndex < 0 || buttonIndex >= len(buttonWirings) {
		log.Fatalf("Invalid button index: %d", buttonIndex)
	}
	wiring := buttonWirings[buttonIndex]
	for i, toggle := range wiring {
		if toggle {
			indicatorLights[i] = !indicatorLights[i]
		}
	}
}

func matchesPattern(indicatorLights []bool, pattern []bool) bool {
	if len(indicatorLights) != len(pattern) {
		return false
	}
	for i := range indicatorLights {
		if indicatorLights[i] != pattern[i] {
			return false
		}
	}
	return true
}

func parseWirings(wirings [][]int, numLights int) [][]bool {
	result := make([][]bool, len(wirings))
	for i := range result {
		result[i] = make([]bool, numLights)
	}
	for i, wiring := range wirings {
		for _, lightIndex := range wiring {
			result[i][lightIndex] = true
		}
	}
	return result
}

func findBestButtonPresses(machine MachineConfig) int {
	indicatorLights := make([]bool, len(machine.indicatorLights))
	targetPattern := machine.indicatorLights
	buttonWirings := machine.buttonWirings
	// given an input config of indicator lights and button wirings,
	// find the minimum number of button presses to match the target pattern
	// use bitwise combinations to try all button press combinations
	// XOR the button wirings to the indicator lights
	// track the minimum number of presses that result in a match
	minPresses := -1
	numButtons := len(buttonWirings)
	totalCombinations := 1 << numButtons // 2^numButtons

	for combo := 0; combo < totalCombinations; combo++ {
		// reset lights to initial state
		currentLights := make([]bool, len(indicatorLights))
		copy(currentLights, indicatorLights)
		pressCount := 0

		for buttonIndex := 0; buttonIndex < numButtons; buttonIndex++ {
			if (combo & (1 << buttonIndex)) != 0 {
				pressButton(buttonIndex, currentLights, buttonWirings)
				pressCount++
			}
		}

		if matchesPattern(currentLights, targetPattern) {
			log.Printf("Found matching combination: %b with %d presses\n", combo, pressCount)
			if minPresses == -1 || pressCount < minPresses {
				minPresses = pressCount
			}
		}
	}

	if minPresses != -1 {
		log.Printf("Minimum button presses to match pattern: %d\n", minPresses)
	} else {
		log.Println("No combination of button presses can match the target pattern.")
	}

	return minPresses
}

func findJoltageMatches(machine MachineConfig) int {
	joltages := machine.joltages
	buttonWirings := machine.buttonWirings

	// Linear Programming Problem:
	// Minimize: c^T * x where c = [1, 1, 1, ...] (sum of button presses)
	// Subject to: A * x = b (joltage constraints)
	//             x >= 0 (non-negative presses)
	// where A[i,j] = 1 if button j affects counter i, 0 otherwise
	//       b[i] = target joltage for counter i

	numCounters := len(joltages)
	numButtons := len(buttonWirings)

	log.Printf("Solving LP for %d counters with %d buttons\n", numCounters, numButtons)

	// Create coefficient matrix A from button wirings
	// Each row represents a counter, each column represents a button
	aData := make([]float64, numCounters*numButtons)
	for i := 0; i < numCounters; i++ {
		for j := 0; j < numButtons; j++ {
			if i < len(buttonWirings[j]) && buttonWirings[j][i] {
				aData[i*numButtons+j] = 1.0
			}
		}
	}

	A := mat.NewDense(numCounters, numButtons, aData)

	// Create target joltages vector b
	bData := make([]float64, numCounters)
	for i, jolt := range joltages {
		bData[i] = float64(jolt)
	}

	// Create objective function: minimize sum of x (c = [1, 1, 1, ...])
	c := make([]float64, numButtons)
	for i := range c {
		c[i] = 1.0
	}

	// Print matrix and target vector for debugging
	log.Println("Coefficient Matrix A (button wirings):")
	fmt.Println(mat.Formatted(A))
	log.Printf("Target Joltages Vector b: %v\n", bData)
	log.Printf("Objective function c (all 1s): %v\n", c)

	// Solve using simplex method in standard form
	// Simplex(c, A, b, tol, initialBasic)
	// Minimizes: c^T * x
	// Subject to: A*x = b, x >= 0
	// Note: A must have full row rank (# rows <= # cols)
	tolerance := 1e-10
	var initialBasic []int // nil means auto-detect

	solutionF, solutionX, err := lp.Simplex(c, A, bData, tolerance, initialBasic)
	if err != nil {
		log.Printf("Error solving LP: %v\n", err)
		return 0
	}

	log.Println("Solution (button presses):")
	for i, val := range solutionX {
		log.Printf("Button %d: %.6f\n", i, val)
	}
	log.Printf("Optimal value (total presses): %.6f\n", solutionF)

	// Sum up total button presses and round to nearest integer
	totalPresses := 0
	for i, val := range solutionX {
		presses := int(val + 0.5) // Round to nearest integer
		if presses < 0 {
			presses = 0 // Ensure non-negative
		}
		totalPresses += presses
		log.Printf("Button %d: %.2f presses -> %d presses\n", i, val, presses)
	}

	log.Printf("Total button presses needed: %d\n", totalPresses)
	return totalPresses
}
