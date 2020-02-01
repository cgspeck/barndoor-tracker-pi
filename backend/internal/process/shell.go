package process

import (
	"bytes"
	"log"
	"os/exec"
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
			"Exit status %v executing %q:\nCaptured StdOut:%v\nCaptured StdErr%v\n",
			exitError.ExitCode,
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
