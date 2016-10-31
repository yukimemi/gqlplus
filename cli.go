package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"gopkg.in/pipe.v2"
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

	// Get files.
	fs, e := getFiles(".")
	failOnError(e)
	for f := range fs {
		fmt.Println(f)
	}

	// Execute bash.
	cmd, e := cmdLoop("bash")

	// Wait cmd exit.
	failOnError(cmd.Wait())

	_ = sql

	_ = user

	_ = pass

	_ = sid

	return ExitCodeOK
}

func cmdLoop(command string) (*exec.Cmd, error) {
	var e error

	cmd := exec.Command(command)
	stdin, e := cmd.StdinPipe()
	stdout, e := cmd.StdoutPipe()
	stderr, e := cmd.StderrPipe()

	e = cmd.Start()
	if e != nil {
		return nil, e
	}

	// stdout scan loop.
	outScanner := bufio.NewScanner(stdout)
	go scanLoop(outScanner)
	// stderr scan loop.
	errScanner := bufio.NewScanner(stderr)
	go scanLoop(errScanner)

	// Get os.stdin and put cmd stdin.
	go func() {
		p := pipe.Line(
			pipe.Read(os.Stdin),
			pipe.Write(stdin),
		)
		failOnError(pipe.Run(p))
		s, e := pipe.Output(p)
		failOnError(e)
		fmt.Println(s)
	}()

	// return command
	return cmd, e
}

func getFiles(root string) (chan string, error) {
	var (
		q  = make(chan string)
		e  error
		wg = new(sync.WaitGroup)
		fn func(p string)
	)

	fmt.Printf("root: [%s]", root)

	// Get file list func.
	fn = func(p string) {
		defer func() { wg.Done() }()

		fis, e := ioutil.ReadDir(p)
		if e != nil {
			failOnError(e)
		}
		for _, fi := range fis {
			full := filepath.Join(p, fi.Name())
			if fi.IsDir() {
				wg.Add(1)
				go fn(full)
			} else {
				q <- full
			}
		}
	}

	// Start file list.
	wg.Add(1)
	go fn(root)

	// Wait.
	go func() {
		wg.Wait()
		close(q)
	}()

	return q, e
}

func scanLoop(scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if e := scanner.Err(); e != nil {
		fmt.Fprintln(os.Stderr, e)
	}
}

// failOnError is easy to judge error.
func failOnError(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}
