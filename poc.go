package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

func main() {
	sshConf := &ssh.ServerConfig{
		PublicKeyCallback: func(_ ssh.ConnMetadata, _ ssh.PublicKey) (*ssh.Permissions, error) {
			return nil, nil
		},
	}

	listener, err := net.Listen("tcp", "0.0.0.0:2022")
	if err != nil {
		log.Fatal("failed to listen for connection: ", err)
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

	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Fatal("failed to accept incoming connection: ", err)
		}

		_, chans, reqs, err := ssh.NewServerConn(nConn, sshConf)
		if err != nil {
			if strings.Contains(err.Error(), "no auth passed yet") {
				return
			}
			if errors.Is(err, io.EOF) {
				continue
			}
			log.Fatal("failed to handshake: ", err)
		}

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
			os.Exit(0)
		}
	}
}
