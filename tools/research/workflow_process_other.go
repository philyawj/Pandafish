//go:build !darwin && !linux

package main

import (
	"os"
	"os/exec"
	"time"
)

func lockWorkflowFile(path string) (func(), error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	return func() { f.Close(); os.Remove(path) }, nil
}

// Other hosts retain CommandContext's direct-process cancellation. Native engine
// pipeline timeouts require Darwin/Linux process-group cleanup; see the guide.
func configureWorkflowProcess(cmd *exec.Cmd) { cmd.WaitDelay = 5 * time.Second }
