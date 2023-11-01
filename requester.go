package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func clearConsole() {
	// Clear the console window
	cmd := exec.Command("cmd", "/c", "cls") // For Linux, use exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func fetchHeartbeat() {
	resp, err := http.Get("http://localhost:7777/heartbeat")
	if err != nil {
		clearConsole()
		fmt.Printf("Error fetching heartbeat: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		clearConsole()
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	clearConsole()
	fmt.Println(string(body))
}

func main() {
	// Set up a ticker to tick every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fetchHeartbeat()
	}
}
