package main

import (
	"bytes"
	"strings"
	"testing"
)

func newTestCLI() (*CLI, *bytes.Buffer, *bytes.Buffer) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	terminate := func(status int) {}
	cli := &CLI{outStream: outStream, errStream: errStream, terminate: terminate}
	return cli, outStream, errStream
}

func prout(command string) []string {
	return append([]string{"git-prout"}, strings.Split(command, " ")...)
}

// FIXME: fetch PR tests...

func TestCLIRun_InvalidRemote(t *testing.T) {
	reset := createTestRepo()
	defer reset()

	cli, _, errStream := newTestCLI()
	cli.Run(prout("--remote foo 10"))

	if got, want := errStream.String(), "'foo' is invalid"; !strings.Contains(got, want) {
		t.Fatalf("Run output %q want %q", got, want)
	}
}

func TestCLIRun_VersionFlag(t *testing.T) {
	cli, _, errStream := newTestCLI()
	cli.Run(prout("--version"))

	if got, want := errStream.String(), Version; got == want {
		t.Fatalf("Run output %q want %q", got, want)
	}
}
