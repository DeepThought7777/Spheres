package tricore

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"spheres/toolkit"
)

type TriCore struct {
	SetName string   `json:"setname"`
	Names   []string `json:"names"`
	Index   int      `json:"index"`
}

// NewTriCore builds a new TriCore object, based on the passed in JSON filename
func NewTriCore(filePath string, index int) (*TriCore, error) {
	data, err := toolkit.ReadFile(filePath)
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

		time.Sleep(toolkit.CheckTime * time.Millisecond)
		err = t.KeepOthersAlive()
	}
}

// WriteLifeSign writes the life sign file for the current core, in UnixMilli style
func (t *TriCore) WriteLifeSign(index int) error {
	filename := t.getLastSeenFilename(index)

	_, err := toolkit.WriteFile(filename, []byte(toolkit.GetUnixMilli()))
	if err != nil {
		return fmt.Errorf("error writing timestamp file")
	}
	return nil
}

// KeepOthersAlive goes over all cores, checking each of its peers and restarting them if needed.
func (t *TriCore) KeepOthersAlive() error {
	for i, _ := range t.Names {
		if i != t.Index {
			_ = t.CheckAndOptionallyStart(i)
		}
	}
	return nil
}

// CheckAndOptionallyStart checks the health of the given core, and restarts it if needed
func (t *TriCore) CheckAndOptionallyStart(index int) error {
	if !t.CheckNodeHealth(index) {
		return t.startPeer(index)
	}
	return nil
}

// startPeer starts the core of the TriCore corresponding to the index
func (t *TriCore) startPeer(index int) error {
	cmd, err := t.runPlatformSpecificScript(index)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("error starting spheres program: %v", err)
	}

	return nil
}

// CheckNodeHealth checks if the given core is still alive
func (t *TriCore) CheckNodeHealth(index int) bool {
	content, err := toolkit.ReadFile(t.getLastSeenFilename(index))
	if err != nil {
		return false
	}

	lastSeenTime := toolkit.StringToUnixMilli(string(content))
	if lastSeenTime == -1 {
		return false
	}

	nowTime := time.Now().UnixMilli()
	fmt.Printf("%s %d\n", t.Names[index], toolkit.DifferenceBetweenUnixMilliTimes(lastSeenTime, nowTime))
	if toolkit.DifferenceBetweenUnixMilliTimes(lastSeenTime, nowTime) > toolkit.CheckTime*2 {
		err = t.WriteLifeSign(index)
		if err != nil {
			fmt.Printf("could not reset timestamp for %s\n", t.Names[index])
		}

		return false
	}
	return true
}

// getLastSeenFilename builds the filename for the lastseen file the core at the given index uses
func (t *TriCore) getLastSeenFilename(index int) string {
	setName := strings.Replace(t.SetName, ".json", "", 1)
	return "LastSeen" + setName + t.Names[index] + ".txt"
}

// runPlatformSpecificScript runs an OS command script to start a new spheres instance
func (t *TriCore) runPlatformSpecificScript(index int) (*exec.Cmd, error) {
	switch runtime.GOOS {
	case "windows":
		// Windows platform
		return exec.Command(".\\startsphere.cmd", t.Names[index], t.SetName, fmt.Sprintf("%d", index)), nil
	case "linux":
		// Linux platform
		return exec.Command("lxterminal", "-e", ".\\sphere", t.SetName, fmt.Sprintf("%d", index)), nil
	case "darwin":
		// macOS (Apple) platform
		return exec.Command(".\\startsphere.sh", t.Names[index], t.SetName, fmt.Sprintf("%d", index)), nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
