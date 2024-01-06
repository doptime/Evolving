package main

import (
	"os"
	"os/exec"
)

// the main function receives a  command ,and it's command line arguments, and run the command with the arguments
func main() {
	// the first argument is the command
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	// set the standard input, output, and error to the current process's standard input, output, and error
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// run the command
	cmd.Run()
}
