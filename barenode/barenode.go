package barenode

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"spheres/toolkit"
	"spheres/tricore"
	"time"
)

const CheckTime = 3

type BareNode struct {
	tricore.TriCore `json:"tricore"`
	NodeGuid        string `json:"nodeguid"`
	ServerPort      int    `json:"serverport"`
}

type HeartbeatResponse struct {
	GUID   string `json:"guid"`
	Server string `json:"server"`
}

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

		time.Sleep(CheckTime * time.Second)
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

func (b *BareNode) Server() {
	mux := http.NewServeMux() // Create a new ServeMux

	mux.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		response := HeartbeatResponse{
			GUID:   b.NodeGuid,
			Server: b.Names[b.Index],
		}

		w.Header().Set("Content-Type", "application/json")

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return
		}

		_, _ = w.Write(jsonResponse)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", b.ServerPort),
		Handler: mux, // Use the new ServeMux as the server handler
	}

	if err := server.ListenAndServe(); err != nil {
		return
	}
}
