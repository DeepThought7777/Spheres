package barenode

import (
	"github.com/google/uuid"
	"os"
	"strings"
)

type BareNode struct {
	// fields and methods specific to BareNode
}

func getProgramGUID(filename string) (string, error) {
	guidBytes, err := os.ReadFile(filename)
	if err != nil {
		return createProgramGUID(filename)
	}

	existingGUID := strings.TrimSpace(string(guidBytes))
	_, err = uuid.Parse(existingGUID)
	if err != nil {
		return createProgramGUID(filename)
	}

	return existingGUID, nil
}

func createProgramGUID(filename string) (string, error) {
	newGUID := uuid.New().String()
	err := os.WriteFile(filename, []byte(newGUID), 0644)
	if err != nil {
		return "", err
	}
	return newGUID, nil
}
