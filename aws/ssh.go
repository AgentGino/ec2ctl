package aws

import (
	"fmt"
	"os"
	"os/exec"

	logger "ec2ctl/logger"

	"github.com/spf13/cobra"
)

func SSHInstance(cmd *cobra.Command, args []string) {
	sshConfigPath := fmt.Sprintf("%s.config", args[0])

	// Check if config file exists
	if _, err := os.Stat(sshConfigPath); os.IsNotExist(err) {
		logger.Error("SSH config file %s does not exist", sshConfigPath)
	}

	// Print the SSH command that will be executed
	fmt.Printf("Executing SSH command with config from: %s\n", sshConfigPath)

	// Add more verbose flags for debugging
	sshCmd := exec.Command("ssh",
		"-v", // More verbose output
		"-F", sshConfigPath,
		"-o", "StrictHostKeyChecking=no", // Disable host key checking
		"-o", "ConnectTimeout=10", // Set connection timeout
		args[0])

	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr
	sshCmd.Stdin = os.Stdin // Add this to allow interactive shell

	err := sshCmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			logger.Error("SSH command failed with exit code %d: %v", exitError.ExitCode(), err)
		}
		logger.Error("failed to execute SSH command: %v", err)
	}
}
