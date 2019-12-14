package omssh

import (
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

// Device : device interface
type Device interface {
	SSHConnect(config *ssh.ClientConfig) error
	SetupIO()
	StartShell() error
	Close() error
}

// SSHDevice : ssh device
type SSHDevice struct {
	Host    string
	Port    string
	client  *ssh.Client
	session *ssh.Session
}

// NewSSHDevice : new SSH device
func NewSSHDevice(host, port string) Device {
	return &SSHDevice{
		Host: host,
		Port: port,
	}
}

// ConfigureSSHClient : configure ssh client
func ConfigureSSHClient(user string, signer ssh.Signer) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

// SSHConnect : ssh connect
func (d *SSHDevice) SSHConnect(config *ssh.ClientConfig) error {
	target := net.JoinHostPort(d.Host, d.Port)
	client, err := ssh.Dial("tcp", target, config)
	if err != nil {
		return err
	}
	d.client = client

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	d.session = session
	return nil
}

// SetupIO : set I/O
func (d *SSHDevice) SetupIO() {
	d.session.Stdout = os.Stdout
	d.session.Stderr = os.Stderr
	d.session.Stdin = os.Stdin
}

// StartShell : start shell
func (d *SSHDevice) StartShell() error {
	defer func() {
		err := d.session.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	term := os.Getenv("TERM")
	err := d.session.RequestPty(term, 25, 80, modes)
	if err != nil {
		return err
	}
	err = d.session.Shell()
	if err != nil {
		return err
	}

	err = d.session.Wait()
	if err != nil {
		return err
	}
	return nil
}

// Close : close client
func (d *SSHDevice) Close() error {
	return d.client.Close()
}
