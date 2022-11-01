package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)

func serveSSH() {
	sshConf := &ssh.ServerConfig{
		PublicKeyCallback: func(_ ssh.ConnMetadata, _ ssh.PublicKey) (*ssh.Permissions, error) {
			return nil, nil
		},
	}

	privateBytes, err := os.ReadFile("id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key: ", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}

	sshConf.AddHostKey(private)

	listener, err := net.Listen("tcp", "localhost:2022")
	if err != nil {
		log.Fatal("failed to listen for connection: ", err)
	}

	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept incoming connection: %v", err)
			continue
		}

		go handleSSHSession(sshConf, nConn)
	}
}

func handleSSHSession(sshConf *ssh.ServerConfig, nConn net.Conn) {
	defer nConn.Close()

	serverConn, chans, reqs, err := ssh.NewServerConn(nConn, sshConf)
	if err != nil {
		if strings.Contains(err.Error(), "no auth passed yet") {
			return
		}
		if errors.Is(err, io.EOF) {
			return
		}
		log.Printf("failed to handshake: %v", err)
		return
	}

	defer serverConn.Close()

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Fatalf("Could not accept channel: %v", err)
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				req.Reply(req.Type == "shell", nil)
			}
		}(requests)

		fmt.Fprintln(channel, "Hello, world!")
		channel.Close()
	}
}

func main() {
	go serveSSH()

	cmd := exec.Command("ssh", "-V")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("could not check SSH version; not in PATH?: %v", err)
	}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		log.Print(line)
	}

	cmd = exec.Command("ssh",
		"-F", "/dev/null",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "StrictHostKeyChecking=no",
		"-o", "IdentityAgent=/dev/null",
		"-o", "IdentityFile=id_rsa_client",
		"-p", "2022",
		"-v",
		"localhost")

	output, err = cmd.CombinedOutput()
	if strings.Contains(string(output), "no mutual signature algorithm") {
		log.Print("connection failed with RSA private key")
		for _, line := range strings.Split(string(output), "\n") {
			if strings.Contains(line, "no mutual signature algorithm") {
				log.Print(line)
			}
		}
	} else if err != nil {
		for _, line := range strings.Split(string(output), "\n") {
			log.Print(line)
		}
		log.Fatal(err)
	}

	cmd = exec.Command("ssh",
		"-F", "/dev/null",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "StrictHostKeyChecking=no",
		"-o", "IdentityAgent=/dev/null",
		"-o", "IdentityFile=id_ed25519_client",
		"-p", "2022",
		"-v", "localhost")

	output, err = cmd.CombinedOutput()
	if strings.Contains(string(output), "Hello, world!") {
		log.Print("connection succeeded with ed25519 private key")
	} else {
		for _, line := range strings.Split(string(output), "\n") {
			log.Print(line)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}
