package main

import (
	"fmt"
	"os"
	"spheres/barenode"
	"spheres/toolkit"
	"strconv"
)

// main() reads the command line arguments and starts the specified member of the TriCore set,
// which is essentially a set of three spheres.go programs guarding each other's backs.
func main() {
	coreFilename, coreIndex, err := getCommandLineArguments()
	if err != nil {
		toolkit.DisplayAndOptionallyExit(err.Error(), true)
	}

	bareNode, err := barenode.NewBareNode(coreFilename, coreIndex)
	if err != nil {
		toolkit.DisplayAndOptionallyExit("TriCore could not be instantiated: "+err.Error(), true)
	}

	bareNode.Run(coreIndex)
}

// getCommandLineArguments() retrieves the command line arguments for main(), and formats them correctly.
// Returns:
// - a string containing the filename of the JSON names set containing the three node names.
// - an integer determining which of the three nodes to start
// - an error if not able to get the proper command line arguments
func getCommandLineArguments() (string, int, error) {
	if len(os.Args) != 3 {
		return "", 0, fmt.Errorf("usage: go run spheres.go <json_file_name> <0-2>")
	}

	jsonFilename := os.Args[1]
	numStr, err := ConvertAndValidateRange(os.Args[2], 0, 2)
	if err != nil {
		return "", 0, fmt.Errorf("usage: go run spheres.go <json_file_name> <0-2>")
	}

	return jsonFilename, numStr, nil
}

func ConvertAndValidateRange(number string, min, max int) (int, error) {
	num, err := strconv.Atoi(number)
	if err != nil {
		return 0, fmt.Errorf("Invalid number provided. Please enter a number between %d and %d.", min, max)
	}

	if num < min || num > max {
		return 0, fmt.Errorf("Invalid number provided. Please enter a number between %d and %d.", min, max)
	}
	return num, nil
}
