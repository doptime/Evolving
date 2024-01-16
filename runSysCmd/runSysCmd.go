package runSysCmd

import (
	"os"
	"os/exec"
	"strings"
)

// the RunCmd function receives a  command ,and it's command line arguments, and run the command with the arguments
func RunCmd(command string) string {
	if len(command) <= 1 {
		return ""
	}
	// the first argument is the command
	cmd := exec.Command("sh", "-c", command)
	// set the standard input, output, and error to the current process's standard input, output, and error
	cmd.Stdin = os.Stdin
	//create memory stdOut with 64k bytes
	stdOut := &strings.Builder{}
	cmd.Stdout = stdOut // Capture the command's output and store it in stdOut
	cmd.Stderr = os.Stderr
	// run the command
	cmd.Run()
	result := stdOut.String()
	return result
}
