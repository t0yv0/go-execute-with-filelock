package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/rogpeppe/go-internal/lockedfile"
)

func main() {
	lockFile := flag.String("lockfile", "", "path to lock file")
	flag.Parse()
	args := flag.Args()

	if *lockFile == "" {
		log.Fatal(fmt.Errorf("-lockedfile cannot be empty"))
	}

	exitCode := executeLocked(*lockFile, args[0], args[1:])
	os.Exit(exitCode)
}

func executeLocked(lockFile string, cmdName string, args []string) int {
	mutex := lockedfile.MutexAt(lockFile)
	unlock, err := mutex.Lock()
	if err != nil {
		log.Fatal(err)
	}
	defer unlock()

	cmd := exec.Command(cmdName, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
	return cmd.ProcessState.ExitCode()
}
