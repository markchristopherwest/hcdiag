package seeker

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"github.com/kballard/go-shellquote"
)

// Commander runs shell commands.
type Commander struct {
	Command string `json:"command"`
	format  string
}

// NewCommander provides a Seeker for running shell commands.
func NewCommander(command string, format string) *Seeker {
	return &Seeker{
		Identifier: command,
		Runner:     Commander{Command: command, format: format},
	}
}

func IsCommandAvailable(name string) bool {
	path, err := exec.LookPath(name)
	if err != nil {
		fmt.Printf("exec.LookPath error: %s", err)
		return false
	}

	log.Println(path) // bin/ls

	return true
}

func (c Commander) Run() (result interface{}, err error) {
	bits, err := shellquote.Split(c.Command)
	if err != nil {
		return nil, fmt.Errorf("Failed to split command: %w", err)
	}

	cmd := bits[0]
	args := bits[1:]

	bts, err := exec.Command(cmd, args...).CombinedOutput()

	if err != nil {
		err = fmt.Errorf("exec.Command error: %s", err)
	}

	switch {
	case c.format == "string":
		result = strings.TrimSuffix(string(bts), "\n")

	case c.format == "json":
		if e := json.Unmarshal(bts, &result); e != nil {
			err = fmt.Errorf("json.Unmarshal error: %s", e)
		}

	default:
		err = errors.New("command output format must be either 'string' or 'json'")
	}

	return result, err
}
