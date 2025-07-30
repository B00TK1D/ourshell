package main

import (
	"github.com/creack/pty"
	"io"
	"log"
	"net"
	"os/exec"
)

func main() {
	// Listen on TCP port 1337
	listener, err := net.Listen("tcp", "0.0.0.0:1337")
	if err != nil {
		log.Fatalf("Failed to listen on port 1337: %v", err)
	}
	defer listener.Close()
	log.Println("Listening on port 1337...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// Handle each connection in a separate goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Start the tmux process with a PTY
	cmd := exec.Command("/usr/local/bin/tmux")

	// Create a pseudo-terminal
	ptyFile, err := pty.Start(cmd)
	if err != nil {
		log.Printf("Failed to start PTY: %v", err)
		return
	}
	defer ptyFile.Close()

	// Forward data between the connection and the PTY
	go func() {
		_, _ = io.Copy(ptyFile, conn)
	}()
	_, _ = io.Copy(conn, ptyFile)

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		log.Printf("Command exited with error: %v", err)
	}
}
