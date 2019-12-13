package omssh

import (
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

// Device : represents a remote network device.
type Device struct {
	Host    string
	Port    string
	client  *ssh.Client
	session *ssh.Session
}

// ConfigureSSHClient : configure ssh client
func ConfigureSSHClient(user string, signer ssh.Signer) *ssh.ClientConfig {
	auth := []ssh.AuthMethod{}
	auth = append(auth, ssh.PublicKeys(signer))

	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return sshConfig
}

func (d *Device) SSHConnect(config *ssh.ClientConfig) error {
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

func (d *Device) SetupIO() {
	d.session.Stdout = os.Stdout
	d.session.Stderr = os.Stderr
	d.session.Stdin = os.Stdin
}

func (d *Device) StartShell() error {
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
