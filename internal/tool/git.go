package tool

import (
	"os/exec"
	"strings"
)

func StagedDiff() (string, error) {
	out, err := exec.Command("git", "diff", "--staged").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func CommitLog() (string, error) {
	// Try to get log from merge base with main
	out, err := exec.Command("bash", "-c",
		`git log --oneline $(git merge-base HEAD main)..HEAD 2>/dev/null || git log --oneline -10`).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
