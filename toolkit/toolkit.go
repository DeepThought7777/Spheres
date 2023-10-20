package toolkit

import (
	"fmt"
	"os"
	"strconv"
)

// displayAndOptionallyExit() displays an error message, waits for enter to be pressed, then exits the program
func DisplayAndOptionallyExit(errorMessage string, exit bool) {
	fmt.Println(errorMessage)
	fmt.Println(">>> Press the [ENTER] key to end the program <<<")
	_, err := fmt.Scanln()
	if !exit || err != nil {
		return
	}
	os.Exit(-1)
}

func ConvertAndValidateRange(number string, min, max int) (int, error) {
	num, err := strconv.Atoi(number)
	if err != nil {
		return 0, fmt.Errorf("Invalid number provided. Please enter a number between 1 and 3.")
	}

	if num < min || num > max {
		return 0, fmt.Errorf("Invalid number provided. Please enter a number between %d and %d.", min, max)
	}
	return num, nil
}

func SecondsBetweenUnixTimes(tThen, tNow int) int {
	return int(tNow - tThen)
}
