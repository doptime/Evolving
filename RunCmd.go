package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// the RunCmd function receives a  command ,and it's command line arguments, and run the command with the arguments
func runCmd(command string, args ...string) string {
	// the first argument is the command
	cmd := exec.Command(command, args...)
	// set the standard input, output, and error to the current process's standard input, output, and error
	cmd.Stdin = os.Stdin
	//create memory stdOut
	var stdOut = &bytes.Buffer{}
	cmd.Stdout = stdOut
	cmd.Stderr = os.Stderr
	// run the command
	cmd.Run()
	return stdOut.String()
}
func RunJSCmd(js string) (prompt string, err error) {
	type Command struct {
		Command string `json:"cmd"`
	}
	var (
		cmd            *Command
		commandStrings []string
	)

	if err = json.Unmarshal([]byte(js), &cmd); err != nil {
		return "", err
	}

	//print the question to the console, and wait for the answer
	fmt.Println(" the product manager has a question: " + cmd.Command)
	//read the answer from the console

	if commandStrings = strings.Split(cmd.Command, " "); len(commandStrings) == 0 {
		fmt.Println("no command")
		return "", nil
	}

	return runCmd(commandStrings[0], commandStrings[1:]...), nil
}
