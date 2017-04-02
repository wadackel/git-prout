package main

import "testing"

func TestPR(t *testing.T) {
	pr := NewPR("hoge", 10, false)

	if want := "hoge"; pr.Remote != want {
		t.Fatalf("pr.Remote got %q want %q", pr.Remote, want)
	}

	if want := 10; pr.Number != want {
		t.Fatalf("pr.Number got %q want %q", pr.Number, want)
	}

	if want := "pr/10"; pr.Branch != want {
		t.Fatalf("pr.Branch got %q want %q", pr.Branch, want)
	}

	if want := "pull/10/head:pr/10"; pr.Ref != want {
		t.Fatalf("pr.Ref got %q want %q", pr.Ref, want)
	}
}
