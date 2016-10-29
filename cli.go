package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object.
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
	cmdFile              string
	cmdDir               string
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		sql  string
		user string
		pass string
		sid  string

		version bool

		e error
	)

	// Get cmd info.
	cli.cmdFile, e = filepath.Abs(os.Args[0])
	failOnError(e)
	cli.cmdDir = filepath.Dir(cli.cmdFile)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&sql, "q", "", "input sql file.")
	flags.StringVar(&user, "u", "", "user name.")
	flags.StringVar(&pass, "p", "", "password.")
	flags.StringVar(&sid, "s", "", "database SID.")

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

	// Check arg
	if flags.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [option] [file]\n", os.Args[0])
		flags.PrintDefaults()
		fmt.Fprintf(os.Stderr, "  file:\n        Spooled file.")
		return ExitCodeError
	}

	fmt.Printf("Cmd file: [%s]\n", cli.cmdFile)
	fmt.Printf("Cmd dir : [%s]\n", cli.cmdDir)

	cmd := exec.Command("bash", "-c")
	stdin, e := cmd.StdinPipe()
	stdout, e := cmd.StdoutPipe()
	stderr, e := cmd.StderrPipe()

	_ = stdout
	_ = stderr

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		l := scanner.Text()
		fmt.Println(l)
		io.WriteString(stdin, l)
		stdin.Close()
	}

	_ = sql

	_ = user

	_ = pass

	_ = sid

	return ExitCodeOK
}

// failOnError is easy to judge error.
func failOnError(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}
