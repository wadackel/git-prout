package main

import "os"

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr, terminate: os.Exit}
	cli.Run(os.Args)
}
