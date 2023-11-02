package main

import (
	"fmt"
	"io"
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

func fetchHeartbeat(port string) {
	resp, err := http.Get("http://localhost:" + port + "/heartbeat")
	if err != nil {
		clearConsole()
		fmt.Printf("Error fetching heartbeat: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		clearConsole()
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	clearConsole()
	fmt.Println(string(body))
}

func main() {
	port := "7777"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	fmt.Printf("Checking port: %s\n", port)
	time.Sleep(3 * time.Second)
	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fetchHeartbeat(port)
	}
}
