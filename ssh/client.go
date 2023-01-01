package ssh

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

const TCPTimeout = time.Second * 30

var ErrInvalidDest = errors.New("invalid destination")

type Client struct {
	*ssh.Client
}

func NewClient(dst string, auth []ssh.AuthMethod) (*Client, error) {
	user, addr, err := parseDest(dst)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dst %s: %w", dst, err)
	}

	c, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            user,
		Timeout:         TCPTimeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            auth,
	})
	if err != nil {
		return nil, err
	}

	return &Client{c}, nil
}

func parseDest(dst string) (user, addr string, err error) {
	split := strings.Split(dst, ":")
	if len(split) > 2 {
		err = ErrInvalidDest
		return
	}

	port := "22"
	if len(split) == 2 {
		port = split[1]
	}
	dst = split[0] + ":" + port

	split = strings.Split(dst, "@")
	if len(split) > 2 {
		err = ErrInvalidDest
		return
	}

	user = "root"
	addr = split[0]

	if len(split) == 2 {
		user = split[0]
		addr = split[1]
	}

	err = nil
	return
}
