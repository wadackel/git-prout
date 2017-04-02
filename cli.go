package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	emoji "gopkg.in/kyokomi/emoji.v1"
)

// Status code
const (
	ExitCodeOK = iota
	ExitCodeParseFlagsError
	ExitCodeError
	ExitCodeNotFoundGit
	ExitCodeOutsideWorkTree
	ExitCodeInvalidRemote
	ExitCodeFailedFetch
	ExitCodeFailedUpdate
	ExitCodeFailedCheckout
)

// Kingpin app, and flags and args.
var (
	app = kingpin.New("git-prout", "").Version(Version)

	debug  = app.Flag("debug", "Enable debug mode.").Bool()
	remote = app.Flag("remote", "Reference of remote.").Short('r').HintAction(GitListRemotes).Default("origin").String()
	force  = app.Flag("force", "Force execute pull or checkout.").Short('f').Bool()
	quiet  = app.Flag("quiet", "Silence any progress and errors (other than parse error).").Short('q').Bool()
	number = app.Arg("number", "ID number of pull request").Required().Int()
)

// Create a new spinner.
func newSpinner(w io.Writer, msg string, complete string) *spinner.Spinner {
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgCyan).SprintFunc()

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Writer = w
	s.Color("cyan")
	s.Suffix = " " + msg + "..."
	s.FinalMSG = blue(string(0x2713)) + " " + green(complete) + "\n"

	return s
}

// Dummy writer.
type silentWriter struct{}

func (w *silentWriter) Write([]byte) (int, error) {
	return 0, nil
}

// CLI is command Runner.
type CLI struct {
	outStream io.Writer
	errStream io.Writer
	terminate func(status int)
}

func (cli *CLI) printDebug(msg string) *CLI {
	if *debug {
		log.Printf("%s %s", color.New(color.FgRed).Sprint("[debug]"), msg)
	}
	return cli
}

func (cli *CLI) printError(msg string) *CLI {
	color.New(color.FgWhite, color.BgRed).Fprintf(cli.errStream, "\rError: %s\n", msg)
	return cli
}

func (cli *CLI) printSuccess(msg string) *CLI {
	emoji.Fprintf(cli.outStream, "\n:sparkles: %s\n%s\n", color.New(color.FgGreen).Sprint("Done!"), msg)
	return cli
}

// Run processing based on arguments.
func (cli *CLI) Run(args []string) {

	// Initialize
	app.ErrorWriter(cli.errStream)
	app.UsageWriter(cli.errStream)
	app.Terminate(cli.terminate)
	app.HelpFlag.Short('h')
	app.UsageTemplate(helpText)

	// Parse
	if _, err := app.Parse(args[1:]); err != nil {
		cli.printError(err.Error())
		cli.terminate(ExitCodeParseFlagsError)
		return
	}

	cli.printDebug("Enable debug mode.")

	// Silent
	if *quiet {
		cli.outStream = &silentWriter{}
		cli.errStream = &silentWriter{}
		cli.printDebug("Enable quiet mode.")
	}

	// Phase 1: Check
	s := newSpinner(cli.outStream, "Checking", "Checked")
	s.Start()

	if !hasGitCommand() {
		cli.printError("'git' command is required.")
		cli.terminate(ExitCodeNotFoundGit)
		return
	}

	if !isInsideGitWorkTree() {
		cli.printError("'git-prout' needs to be executed in work tree.")
		cli.terminate(ExitCodeOutsideWorkTree)
		return
	}

	if !GitIsValidRemote(*remote) {
		cli.printError(fmt.Sprintf("'%s' is invalid remote.", *remote))
		cli.printDebug("remotes -> [" + strings.Join(GitListRemotes(), ", ") + "]")
		cli.terminate(ExitCodeInvalidRemote)
		return
	}

	currentBranch, err := GitCurrentBranch()
	if err != nil {
		cli.printError("Failed to acquire the current branch.")
		cli.printDebug(err.Error())
		cli.terminate(ExitCodeError)
		return
	}

	pr := NewPR(*remote, *number, *force)  // PR
	isUpdate := pr.Branch == currentBranch // Mode

	s.Stop()
	cli.printDebug(fmt.Sprintf("Cheked, isUpdate = %t", isUpdate))

	// Phase 2: Fetch
	s = newSpinner(cli.outStream, "Fetching", "Fetched")
	s.Start()
	if _, err := pr.Fetch(); err != nil {
		cli.printError(fmt.Sprintf("Failed to fetch remote ref '%s %s'.", pr.Remote, pr.Ref))
		cli.printDebug(err.Error())
		cli.terminate(ExitCodeFailedFetch)
		return
	}
	s.Stop()
	cli.printDebug(fmt.Sprintf("Fetched, %s %s.", pr.Remote, pr.Ref))

	// Phase 3: Update or Checkout
	if isUpdate {
		s = newSpinner(cli.outStream, "Updating", "Updated")
		s.Start()
		if _, err := pr.Apply(); err != nil {
			cli.printError("Failed to update.")
			cli.printDebug(err.Error())
			cli.terminate(ExitCodeFailedUpdate)
			return
		}

	} else {
		s = newSpinner(cli.outStream, "Checkout", "Checkout")
		s.Start()
		if _, err := pr.Checkout(); err != nil {
			cli.printError("Failed to checkout.")
			cli.printDebug(err.Error())
			cli.terminate(ExitCodeFailedCheckout)
			return
		}
	}
	s.Stop()

	// Success
	var msg string

	if isUpdate {
		msg = fmt.Sprintf("Updated a '%s' branch.", pr.Branch)
	} else {
		msg = fmt.Sprintf("Switched to a '%s' branch.", pr.Branch)
	}

	cli.printSuccess(msg)
}

// for Usage (help).
var title = color.New(color.FgYellow).SprintFunc()

var helpText = fmt.Sprintf(`{{define "FormatCommand"}} [options] {{range .Args}}\
{{if not .Required}}[{{end}}<{{.Name}}>{{if .Value|IsCumulative}}...{{end}}{{if not .Required}}]{{end}}{{end}}\
{{end}}\
{{define "FormatCommands"}}\
{{range .FlattenedCommands}}\
{{if not .Hidden}}\
  {{.FullCommand}}{{if .Default}}*{{end}}{{template "FormatCommand" .}}
{{.Help|Wrap 4}}
{{end}}\
{{end}}\
{{end}}\
{{define "FormatUsage"}}\
{{template "FormatCommand" .}}
{{if .Help}}
{{.Help|Wrap 0}}\
{{end}}\
{{end}}\
%s
  {{.App.Name}}{{template "FormatUsage" .App}}
{{if .Context.Flags}}\
%s
{{.Context.Flags|FlagsToTwoColumns|FormatTwoColumns}}
{{end}}\
{{if .Context.Args}}\
%s
{{.Context.Args|ArgsToTwoColumns|FormatTwoColumns}}
{{end}}\
`, title("Usage:"), title("Options:"), title("Arguments:"))
