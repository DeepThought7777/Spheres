package toolkit

import (
	"fmt"
	"os"
)

// DisplayAndOptionallyExit() displays an error message, waits for enter to be pressed, then exits the program
func DisplayAndOptionallyExit(errorMessage string, exit bool) {
	fmt.Println(errorMessage)
	fmt.Println(">>> Press the [ENTER] key to end the program <<<")
	_, err := fmt.Scanln()
	if !exit || err != nil {
		return
	}
	os.Exit(-1)
}

func SecondsBetweenUnixTimes(tThen, tNow int) int {
	return int(tNow - tThen)
}

func WriteFile(filename string, data []byte) (int, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Errorf("File could not be opened: %v", err)
	}

	defer file.Close()

	count, err := file.Write(data)
	if err != nil {
		fmt.Errorf("File could not be written: %v", err)
	}
	return count, nil
}

func ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename) // For read access.
	if err != nil {
		fmt.Errorf("File could not be opened: %v", err)
	}

	defer file.Close()

	data := make([]byte, 2048)
	count, err := file.Read(data)
	if err != nil {
		fmt.Errorf("File could not be read: %v", err)
	}
	return data[:count], nil
}
