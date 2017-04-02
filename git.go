package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

// Check `git` command can be callable.
func hasGitCommand() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// Check is inside work tree in git repository
func isInsideGitWorkTree() bool {
	_, err := git("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false
	}
	return true
}

// Git call `git` command, return standard output as a string.
func git(subcmd string, args ...string) (string, error) {
	gitArgs := append([]string{subcmd}, args...)

	var out bytes.Buffer
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdout = &out
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() != 0 {
				return "", err
			}
		}
	}

	return strings.TrimRight(strings.TrimSpace(out.String()), "\000"), nil
}

// GitRemotes get all remotes.
func GitRemotes() ([]string, error) {
	out, err := git("remote")
	if err != nil {
		return []string{}, err
	}

	return regexp.MustCompile("\r\n|\n\r|\n|\r").Split(out, -1), nil
}

// GitListRemotes get all remotes. In case of error, it returns an empty array.
func GitListRemotes() []string {
	remotes, err := GitRemotes()
	if err != nil {
		return []string{}
	}
	return remotes
}

// GitIsValidRemote check the exists of the specify remote.
func GitIsValidRemote(remote string) bool {
	remotes, err := GitRemotes()
	if err != nil {
		return false
	}

	for _, r := range remotes {
		if r == remote {
			return true
		}
	}

	return false
}

// GitCurrentBranch get current branch name.
func GitCurrentBranch() (string, error) {
	return git("rev-parse", "--abbrev-ref", "HEAD")
}
