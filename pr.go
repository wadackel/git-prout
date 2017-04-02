package main

import "strconv"

// PR pull request
type PR struct {
	Remote string
	Number int
	Ref    string
	Branch string
	Force  bool
}

// NewPR create new PR instance.
func NewPR(remote string, number int, force bool) PR {
	pr := PR{}
	pr.Remote = remote
	pr.Number = number
	pr.Force = force
	pr.Branch = "pr/" + strconv.Itoa(number)
	pr.Ref = "pull/" + strconv.Itoa(number) + "/head:" + pr.Branch
	return pr
}

// Fetch PR. (force fetch)
func (p *PR) Fetch() (string, error) {
	return git("fetch", p.Remote, p.Ref, "-f", "-u")
}

// Checkout to PR branch.
func (p *PR) Checkout() (string, error) {
	args := []string{p.Branch}
	if p.Force {
		args = append(args, []string{"-f"}...)
	}
	return git("checkout", args...)
}

// Apply PR to the working directory.
func (p *PR) Apply() (string, error) {
	return git("reset", "--hard", "HEAD")
}
