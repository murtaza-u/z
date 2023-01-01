package ssh

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type Session struct {
	*ssh.Session
	wg        *sync.WaitGroup
	closehand chan struct{}
}

func (c *Client) NewSession() (*Session, error) {
	sess, err := c.Client.NewSession()
	if err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	term := os.Getenv("TERM")
	if term == "" {
		term = "xterm-256color"
	}

	w, h, err := terminal.GetSize(0)
	if err != nil {
		h = 80
		w = 40
	}

	err = sess.RequestPty(term, h, w, modes)
	if err != nil {
		sess.Close()
		return nil, fmt.Errorf("PTY request failed: %w", err)
	}

	return &Session{
		Session:   sess,
		wg:        new(sync.WaitGroup),
		closehand: make(chan struct{}),
	}, nil
}

func (s *Session) Close() {
	s.closehand <- struct{}{}
	s.Session.Close()
}

func (s *Session) Wait() {
	s.wg.Wait()
}

func (s *Session) WaitAndClose() {
	s.Wait()
	s.Close()
}

func (s *Session) CloseAndWaith() {
	s.Close()
	s.Wait()
}

func (s *Session) SetPipes() error {
	stdin, err := s.StdinPipe()
	if err != nil {
		s.Close()
		return fmt.Errorf("failed to setup stdin: %w", err)
	}

	stdout, err := s.StdoutPipe()
	if err != nil {
		s.Close()
		return fmt.Errorf("failed to setup stdout: %w", err)
	}

	stderr, err := s.StderrPipe()
	if err != nil {
		s.Close()
		return fmt.Errorf("failed to setup stderr: %w", err)
	}

	go io.Copy(stdin, os.Stdin)

	s.wg.Add(2)
	go s.copy(os.Stdout, stdout)
	go s.copy(os.Stderr, stderr)

	return nil
}

var sigmap = map[string]ssh.Signal{
	"aborted":                  "ABRT",
	"alarm clock":              "ALRM",
	"floating point exception": "FPE",
	"hangup":                   "HUP",
	"illegal instruction":      "ILL",
	"interrupt":                "INT",
	"killed":                   "KILL",
	"broken pipe":              "PIPE",
	"quit":                     "QUIT",
	"segmentation fault":       "SEGV",
	"terminated":               "TERM",
	"user defined signal 1":    "USR1",
	"user defined signal 2":    "USR2",
}

func (s *Session) HandleSignals() {
	sigssh := make(chan os.Signal)
	signal.Notify(
		sigssh,
		syscall.SIGABRT, syscall.SIGALRM, syscall.SIGFPE,
		syscall.SIGHUP, syscall.SIGILL, syscall.SIGINT, syscall.SIGKILL,
		syscall.SIGPIPE, syscall.SIGQUIT, syscall.SIGSEGV,
		syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2,
	)

	for {
		select {
		case sig := <-sigssh:
			err := s.Signal(sigmap[sig.String()])
			if err != nil {
				break
			}

		case <-s.closehand:
			break
		}
	}
}

func (s *Session) copy(dst io.Writer, src io.Reader) {
	io.Copy(dst, src)
	s.wg.Done()
}
