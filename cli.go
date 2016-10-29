package main

import (
	"flag"
	"fmt"
	"io"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		file string
		user string
		pass string
		sid  string

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&file, "file", "", "input sql file")
	flags.StringVar(&file, "f", "", "input sql file(Short)")

	flags.StringVar(&user, "user", "", "user name")
	flags.StringVar(&user, "u", "", "user name(Short)")

	flags.StringVar(&pass, "pass", "", "password")
	flags.StringVar(&pass, "p", "", "password(Short)")

	flags.StringVar(&sid, "sid", "", "database SID")
	flags.StringVar(&sid, "s", "", "database SID(Short)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	_ = file

	_ = user

	_ = pass

	_ = sid

	return ExitCodeOK
}
