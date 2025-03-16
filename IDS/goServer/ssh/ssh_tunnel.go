package ssh

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"goServer/config"
	"golang.org/x/crypto/ssh"
)

// SSHTunnel holds the tunnel connection
type Tunnel struct {
	LocalPort int
	sshClient *ssh.Client
	listener  net.Listener
}

// CreateSSHTunnel establishes an SSH tunnel for PostgreSQL
func CreateSSHTunnel() (*Tunnel, error) {
	// Read SSH private key
	key, err := os.ReadFile(config.SSHKey)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %w", err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %w", err)
	}

	// Set up SSH configuration
	sshConfig := &ssh.ClientConfig{
		User:            config.SSHUser,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the SSH server
	sshConn, err := ssh.Dial("tcp", net.JoinHostPort(config.SSHHost, config.SSHPort), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server: %w", err)
	}

	// Start local listener on a random available port
	localListener, err := net.Listen("tcp", "localhost:0") // Bind to any available port
	if err != nil {
		return nil, fmt.Errorf("failed to create local listener: %w", err)
	}

	// Forward local connections to the remote PostgreSQL server
	go func() {
		for {
			localConn, err := localListener.Accept()
			if err != nil {
				log.Println("Failed to accept connection:", err)
				continue
			}

			// Connect to remote PostgreSQL server via SSH
			remoteConn, err := sshConn.Dial("tcp", net.JoinHostPort(config.DBHost, config.DBPort))
			if err != nil {
				log.Println("Failed to connect to remote PostgreSQL:", err)
				localConn.Close()
				continue
			}

			// Pipe data between local and remote connections
			go io.Copy(localConn, remoteConn)
			go io.Copy(remoteConn, localConn)
		}
	}()

	return &Tunnel{
		LocalPort: localListener.Addr().(*net.TCPAddr).Port,
		sshClient: sshConn,
		listener:  localListener,
	}, nil
}

// Close the SSH tunnel
func (t *Tunnel) Close() {
	t.listener.Close()
	t.sshClient.Close()
}
