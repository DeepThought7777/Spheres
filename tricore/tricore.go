package tricore

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"spheres/toolkit"
	"strconv"
	"time"
)

const CheckTime = 3

type TriCore struct {
	SetName string   `json:"setname"`
	Names   []string `json:"names"`
	Index   int      `json:"index"`
}

// NewTriCore builds a new TriCore object, based on the passed in JSON filename
func NewTriCore(filePath string, index int) (*TriCore, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read JSON file: %v", err)
	}

	var core TriCore
	err = json.Unmarshal(data, &core)
	if err != nil {
		return nil, fmt.Errorf("could not process JSON file: %v", err)

	}

	if len(core.Names) != 3 {
		return nil, fmt.Errorf("core file does not have three nodes")
	}

	core.SetName = filePath
	core.Index = index
	return &core, nil
}

// Run is the basic for loop that will have the core write its lifesign,
// and check the life signs of the other cores, resurrecting them if not recent enough.
func (t *TriCore) Run(index int) {
	for {
		err := t.WriteLifeSign(index)
		if err != nil {
			toolkit.DisplayAndOptionallyExit("TriCore could not be instantiated: "+err.Error(), true)
		}

		time.Sleep(CheckTime * time.Second)
		err = t.KeepOthersAlive()
	}
}

// WriteLifeSign writes the life sign file for the current core, in Unix style
func (t *TriCore) WriteLifeSign(index int) error {
	filename := t.getLastSeenFilename(index)
	currentTime := time.Now()

	err := os.WriteFile(filename, []byte(strconv.Itoa(int(currentTime.Unix()))), 644)
	if err != nil {
		return fmt.Errorf("error writing timestamp file")
	}
	return nil
}

// KeepOthersAlive goes over all cores, checking each of its peers and restarting them if needed.
func (t *TriCore) KeepOthersAlive() error {
	for i, _ := range t.Names {
		if i != t.Index {
			t.CheckAndOptionallyStart(t.SetName, i)
		}
	}
	return nil
}

// CheckAndOptionallyStart reads the last life sign,restarting the core if not found or out of date.
func (t *TriCore) CheckAndOptionallyStart(filename string, index int) error {
	content, err := os.ReadFile(t.getLastSeenFilename(index))
	if err != nil {
		return t.startPeer(filename, index)
	}

	lastSeenTime, err := strconv.Atoi(string(content))
	if err != nil {
		return t.startPeer(filename, index)
	}

	nowTime := int(time.Now().Unix())
	fmt.Println(toolkit.SecondsBetweenUnixTimes(lastSeenTime, nowTime))
	if toolkit.SecondsBetweenUnixTimes(lastSeenTime, nowTime) > CheckTime+1 {
		err = t.WriteLifeSign(index)
		if err != nil {
			fmt.Printf("could not reset timestamp for %s\n", t.Names[index])
		}

		return t.startPeer(filename, index)
	}
	return nil
}

func (t *TriCore) startPeer(filename string, index int) error {
	cmd := exec.Command(".\\startsphere.cmd", t.Names[index], filename, fmt.Sprintf("%d", index))

	// Set the appropriate standard input, output, and error streams
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error starting spheres program: %v", err)
	}

	return nil
}

func (t *TriCore) getLastSeenFilename(index int) string {
	return "LastSeen" + t.Names[index] + ".txt"
}
