package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func checkPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func createTestRepo() func() {
	var err error
	tmpdir, err := ioutil.TempDir(".", "git-prout-test")
	checkPanic(err)

	prevdir, err := filepath.Abs(".")
	checkPanic(err)

	os.Chdir(tmpdir)

	_, err = exec.Command("git", "init").Output()
	checkPanic(err)

	_, err = exec.Command("git", "config", "--local", "user.name", "'testuser'").Output()
	checkPanic(err)
	_, err = exec.Command("git", "config", "--local", "user.email", "'testuser@email.com'").Output()
	checkPanic(err)

	tmpfile := "README.md"
	readme, err := os.Create(tmpfile)
	checkPanic(err)
	_, err = readme.WriteString("foo\n")
	checkPanic(err)

	_, err = exec.Command("git", "add", "-A").Output()
	checkPanic(err)
	_, err = exec.Command("git", "commit", "-m", "'testing'").Output()
	checkPanic(err)
	_, err = exec.Command("git", "remote", "add", "origin", "https://github.com/tsuyoshiwada/git-prout.git").Output()
	checkPanic(err)

	return func() {
		os.Chdir(prevdir)
		os.RemoveAll(tmpdir)
	}
}
