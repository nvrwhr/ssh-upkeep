package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// Define a custom flag type to handle multiple -L arguments
type multiFlag []string

func (m *multiFlag) String() string {
	return strings.Join(*m, ", ")
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func main() {
	// Initialize the custom flag type for multiple -L flags
	var portForwards multiFlag
	flag.Var(&portForwards, "L", "Local port forward in format [localPort:remoteHost:remotePort] (e.g., 5001:server.int:5432). Repeatable.")

	// Define SSH user and host argument
	userHost := flag.String("u", "", "SSH user and host in format [user@host] (e.g., ec2-user@bastion.ai)")

	flag.Parse()

	// Ensure required arguments are provided
	if len(portForwards) == 0 || *userHost == "" {
		log.Fatalf("Usage: %s -L [localPort:remoteHost:remotePort] -u [user@host]\n", "ssh-connector")
	}

	// Construct SSH command with multiple -L flags
	sshCommand := "ssh"
	for _, port := range portForwards {
		sshCommand += fmt.Sprintf(" -L %s", port)
	}
	sshCommand += " " + *userHost

	for {
		fmt.Printf("Attempting to connect with SSH command: %s\n", sshCommand)

		// Run the SSH command as a subprocess
		cmd := exec.Command("sh", "-c", sshCommand)

		// Connect command output to standard output
		cmd.Stdout = log.Writer()
		cmd.Stderr = log.Writer()

		// Start the command
		err := cmd.Start()
		if err != nil {
			log.Printf("Failed to start SSH command: %s\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Wait for the command to finish (blocking call)
		err = cmd.Wait()
		if err != nil {
			log.Printf("SSH command interrupted: %s. Reconnecting in 5 seconds...\n", err)
		} else {
			log.Println("SSH command completed successfully. Reconnecting in 5 seconds...")
		}

		time.Sleep(5 * time.Second) // Reconnect delay
	}
}
