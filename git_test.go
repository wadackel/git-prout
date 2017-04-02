package main

import "testing"

func TestGit(t *testing.T) {
	reset := createTestRepo()
	defer reset()

	if files, _ := git("ls-files"); files != "README.md" {
		t.Fatalf("failed git, ls-files -> %q", files)
	}
}

func TestGitRemotes(t *testing.T) {
	reset := createTestRepo()
	defer reset()

	remotes := GitListRemotes()
	if len(remotes) != 1 || remotes[0] != "origin" {
		t.Fatalf("failed get remotes. %d %q", len(remotes), remotes)
	}

	_, err := git("remote", "add", "foo", "http://hoge.com/fuga/baz.git")
	if err != nil {
		panic(err)
	}

	remotes = GitListRemotes()

	if len(remotes) != 2 || remotes[0] != "foo" || remotes[1] != "origin" {
		t.Fatalf("failed get remotes. %d %q", len(remotes), remotes)
	}
}

func TestGitValidRemote(t *testing.T) {
	reset := createTestRepo()
	defer reset()

	if GitIsValidRemote("hoge") {
		t.Fatal("invalid remote of 'hoge'.")
	}

	if !GitIsValidRemote("origin") {
		t.Fatal("invalid remote of 'origin'.")
	}
}

func TestGitCurrentBranch(t *testing.T) {
	reset := createTestRepo()
	defer reset()

	want := "master"
	if got, _ := GitCurrentBranch(); got != want {
		t.Fatalf("current branch got %s want %s", got, want)
	}

	_, err := git("checkout", "-b", "pr/123")
	if err != nil {
		panic(err)
	}

	want = "pr/123"
	if got, _ := GitCurrentBranch(); got != want {
		t.Fatalf("current branch got %s want %s", got, want)
	}
}
