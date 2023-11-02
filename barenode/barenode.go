package barenode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"spheres/toolkit"
	"spheres/tricore"

	"github.com/google/uuid"
)

// BareNode is the struct that encapsulates a TriCore and adds specific data
type BareNode struct {
	tricore.TriCore `json:"tricore"`
	NodeGuid        string `json:"nodeguid"`
	ServerPort      int    `json:"serverport"`
}

// HeartbeatResponse is the JSON format that the heartbeat handler returns as a body
type HeartbeatResponse struct {
	GUID   string `json:"guid"`
	Server string `json:"server"`
}

// NewBareNode creates a new BareNode structure and returns a reference to it
func NewBareNode(filePath string, index int) (*BareNode, error) {
	data, err := toolkit.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read JSON file: %v", err)
	}

	var node BareNode
	err = json.Unmarshal(data, &node)
	if err != nil {
		return nil, fmt.Errorf("could not process JSON file: %v", err)
	}

	err = node.CompleteAndWriteBareNode(filePath, index)
	if err != nil {
		return nil, fmt.Errorf("could not write JSON file: %v", err)
	}
	return &node, nil
}

// Run is the basic for loop that will have the core write its lifesign,
// and check the life signs of the other cores, resurrecting them if not recent enough.
func (b *BareNode) Run(index int) {
	serverRunning := false
	for {
		if !serverRunning {
			serverRunning = true
			go b.Server()
			serverRunning = false
		}

		time.Sleep(toolkit.CheckTime * time.Millisecond)
		err := b.WriteLifeSign(index)
		if err != nil {
			toolkit.DisplayAndOptionallyExit("TriCore could not be instantiated: "+err.Error(), true)
		}

		err = b.KeepOthersAlive()
		if err != nil {
			toolkit.DisplayAndOptionallyExit("TriCore could not be instantiated: "+err.Error(), true)
		}
	}
}

// CompleteAndWriteBareNode fills the remaining fields of the struct and writes out the JSON file
func (b *BareNode) CompleteAndWriteBareNode(filePath string, index int) error {
	b.TriCore.Index = index

	if b.TriCore.SetName == "" {
		b.TriCore.SetName = filePath
	}

	if b.NodeGuid == "" {
		b.NodeGuid = uuid.New().String()
	}

	return b.WriteBareNodeToFile()
}

// WriteBareNodeToFile writes the given BareNode struct to a file as JSON.
func (b *BareNode) WriteBareNodeToFile() error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal node to JSON: %v", err)
	}

	_, err = toolkit.WriteFile(b.TriCore.SetName, data)
	if err != nil {
		return fmt.Errorf("could not write JSON to file: %v", err)
	}

	return nil
}

// Server binds the port and tries to start the server, failing silently if needed.
func (b *BareNode) Server() {
	mux := http.NewServeMux()

	// Use a method expression to convert the heartbeatHandler method into a function
	mux.HandleFunc("/heartbeat", b.heartbeatHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", b.ServerPort),
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		// Handle the error according to your application's requirements.
		// For example, you might log the error or exit the program.
	}
}

// heartbeatHandler is passed to the HandleFunc to handle the heartbeat
func (b *BareNode) heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	response := HeartbeatResponse{
		GUID:   b.NodeGuid,
		Server: b.Names[b.Index],
	}

	w.Header().Set("Content-Type", "application/json")

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		// Handle the error according to your application's requirements.
		// For example, you might log the error or send a different response.
	}
}
