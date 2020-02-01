package process

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func ShellOut(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if exitError, ok := err.(*exec.ExitError); ok {
		log.Printf(
			"Exit status %d executing %q:\nCaptured StdOut:%v\nCaptured StdErr%v\n",
			exitError.ExitCode(),
			command,
			stdout.String(),
			stderr.String(),
		)
		return exitError, stdout.String(), stderr.String()
	}

	if err != nil {
		err = err.(*exec.ExitError)
		log.Printf(
			"%v executing %q:\nCaptured StdOut:%v\nCaptured StdErr%v\n",
			err,
			command,
			stdout.String(),
			stderr.String(),
		)
	}
	return err, stdout.String(), stderr.String()
}

func RunCommands(commands []string) (error, string, string) {
	var builderStdOut strings.Builder
	var builderStdErr strings.Builder

	var err error

	for i, command := range commands {
		builderStdOut.WriteString(fmt.Sprintf("[%v] %q\n", i, command))
		err, stdOut, stdErr := ShellOut(command)
		builderStdOut.WriteString(stdOut)
		builderStdErr.WriteString(stdErr)
		if err != nil {
			break
		}
	}

	return err, builderStdOut.String(), builderStdErr.String()
}
