package main

import "testing"

func TestFilterCommandsSubstring(t *testing.T) {
	all := []string{"shell:Desktop", "shell:Downloads", "shell:Windows"}
	got := filterCommands(all, "down")
	if len(got) != 1 || got[0] != "shell:Downloads" {
		t.Fatalf("filter: got %v", got)
	}
}

func TestFilterKeepsHighlightIndex(t *testing.T) {
	all := sortedShellCommands()
	filtered := filterCommands(all, "app")
	idx := indexOf(filtered, "shell:AppData")
	if idx < 0 {
		t.Fatal("expected AppData in app filter")
	}
	narrower := filterCommands(all, "appdata")
	if indexOf(narrower, "shell:AppData") < 0 {
		t.Fatal("AppData should remain when narrowing filter")
	}
}
