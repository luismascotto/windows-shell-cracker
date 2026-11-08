package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type savedState struct {
	LastCommand string `json:"last_command"`
}

func stateFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "windows-shell-cracker", "state.json"), nil
}

func loadLastCommand() string {
	path, err := stateFilePath()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var st savedState
	if err := json.Unmarshal(data, &st); err != nil {
		return ""
	}
	return st.LastCommand
}

func saveLastCommand(command string) error {
	if command == "" {
		return nil
	}
	path, err := stateFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(savedState{LastCommand: command})
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func indexOf(commands []string, command string) int {
	for i, c := range commands {
		if c == command {
			return i
		}
	}
	return -1
}

func filterCommands(all []string, filter string) []string {
	if filter == "" {
		out := make([]string, len(all))
		copy(out, all)
		return out
	}
	needle := []rune(filter)
	var out []string
	for _, c := range all {
		if containsFold(c, needle) {
			out = append(out, c)
		}
	}
	return out
}

func containsFold(hay string, needle []rune) bool {
	if len(needle) == 0 {
		return true
	}
	hr := []rune(hay)
	for i := 0; i <= len(hr)-len(needle); i++ {
		match := true
		for j, n := range needle {
			if !runeEqualFold(hr[i+j], n) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func runeEqualFold(a, b rune) bool {
	if a == b {
		return true
	}
	if a >= 'A' && a <= 'Z' {
		a += 'a' - 'A'
	}
	if b >= 'A' && b <= 'Z' {
		b += 'a' - 'A'
	}
	return a == b
}

var errNoSelection = errors.New("no selection")
