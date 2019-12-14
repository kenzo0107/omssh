package omssh

import (
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh"
)

const (
	testPort       = "2222"
	testPrivateKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAq6otEnqrpubCsmeTs/xnaayMu6/VtaEsnFLS5qKWR0dpHqORJ0AJ
nWZZooTNciDNHva7NZCl8DfIYvR3YXanpfCcCIwsHf4eJyVTiguUE1QwHmru59Luf+dunu
pq/9QXTo2I/UNgVS0e9M3o3Dpv3/Ymjb1JIth4fTSDz002Y3Ysuw3k6psGT7N73U9u5X6R
9WLTvarqH0NMmU9+Q7NlxW7T9CWrXAbPxTZmEGizanjKtR6LZ3p2pq9SAeRu1IjPvY+aAU
60980qOB7swl4fzz7FWtv5NAHVAFDwY9fmlebuIs+UkeQcR/Ku+iVkVkXU+njFsJnAiEUs
o0Zb0wPYDQAAA9h00TWpdNE1qQAAAAdzc2gtcnNhAAABAQCrqi0Sequm5sKyZ5Oz/GdprI
y7r9W1oSycUtLmopZHR2keo5EnQAmdZlmihM1yIM0e9rs1kKXwN8hi9Hdhdqel8JwIjCwd
/h4nJVOKC5QTVDAeau7n0u5/526e6mr/1BdOjYj9Q2BVLR70zejcOm/f9iaNvUki2Hh9NI
PPTTZjdiy7DeTqmwZPs3vdT27lfpH1YtO9quofQ0yZT35Ds2XFbtP0JatcBs/FNmYQaLNq
eMq1Hotnenamr1IB5G7UiM+9j5oBTrT3zSo4HuzCXh/PPsVa2/k0AdUAUPBj1+aV5u4iz5
SR5BxH8q76JWRWRdT6eMWwmcCIRSyjRlvTA9gNAAAAAwEAAQAAAQBtacLek1dSwqP3p/LJ
difHf9YXTmRNJtRTMqr/m0NjXQ2QHLrIpJU8QF8DKdf0VRnIEYSTCIXrTPKot55bfZAvQO
OCwyzfVPeNBcpwIx8XDsK4sHljQtsGpNCp80mNk3XjeGyG1+nPgDnJ2HAB5jEmMzKxhqLV
1dk+HDmi6FixHTbuB9Y7ej5s42WB6CdazBBo6iT71aQtoBqmPZX0FM+AcnZxafTtEOCdq3
eIqGJgCpqTT1Z1NptIW4lLAOmpKpsIobQUQivxnctpI2yM+/Inn/LlSeF8uSbVFL1jcTEh
5xuI/cZlCEeUiJ7OhB3vJNP2hnjLB8xGipK4F9Mp4yNhAAAAgQCcCQ8jWBFnMF0iyCFyei
LchJ+K3H+hJmT3U4eKVrnFBXRDY+8s2La2DT7Kp4NUvNrMYuWz98+DpoR6VbM1axSDPrvy
6ggjBRbrL8zZ5v6Lq9GcFRx8KuHtXKhYIJxfShh6hl0r3951fSEzIoAL9voORjABgDj63o
MzkkCapPL99AAAAIEA1fxwnwWh6B+4tmTfiKj7dRtVr2VadydHZnC7tHHZ4A7wgE/yHsqn
F8iJwWexTN2znT8RjyiMewSNrESiyFvmgxnCxGg/y8DU0fK4Hmlc6SPw1wx0w0wyVrJRsn
aYqg5DAuAdcAeXe/NV85ACT6Pc0Oe3aevnLDbaIbEk/RzR70kAAACBAM1eiwIN7yKLSQzc
fU0WdoE6G8P7+eNOJrAIxKURd3EgCQTSUSkhXJhUQ7zfYm7HLrvWH8CwuE4efTMLWkVpR7
LBB3/mscp+OsRZrqwr+wj225qRqqCL6E4NZ1WOxsF/HSdA2Wgai5m/3lsKaf+1918jAw/T
OyqzbCA9Ird5ia6lAAAAHGtlbnpvLnRhbmFrYUBwYy0xMi01NjcubG9jYWwBAgMEBQY=
-----END OPENSSH PRIVATE KEY-----`
	testPublicKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCrqi0Sequm5sKyZ5Oz/GdprIy7r9W1oSycUtLmopZHR2keo5EnQAmdZlmihM1yIM0e9rs1kKXwN8hi9Hdhdqel8JwIjCwd/h4nJVOKC5QTVDAeau7n0u5/526e6mr/1BdOjYj9Q2BVLR70zejcOm/f9iaNvUki2Hh9NIPPTTZjdiy7DeTqmwZPs3vdT27lfpH1YtO9quofQ0yZT35Ds2XFbtP0JatcBs/FNmYQaLNqeMq1Hotnenamr1IB5G7UiM+9j5oBTrT3zSo4HuzCXh/PPsVa2/k0AdUAUPBj1+aV5u4iz5SR5BxH8q76JWRWRdT6eMWwmcCIRSyjRlvTA9gN`
)

func TestConfigureSSHClient(t *testing.T) {
	user := "ubuntu"
	signer, err := ssh.ParsePrivateKey([]byte(testPrivateKey))
	if err != nil {
		t.Error("wrong result : err is not nil")
	}

	sshClientConfig := ConfigureSSHClient(user, signer)

	if diff := cmp.Diff(sshClientConfig.User, user); diff != "" {
		t.Errorf("wrong result :\n%s", diff)
	}
}

func TestSSHConnect(t *testing.T) {
	signer, err := ssh.ParsePrivateKey([]byte(testPrivateKey))
	if err != nil {
		t.Error("wrong result : err is not nil")
	}

	go buildSSHServer(signer)

	user := "testUser"
	device := NewDevice("localhost", testPort)
	sshClientConfig := ConfigureSSHClient(user, signer)
	if err := device.SSHConnect(sshClientConfig); err != nil {
		t.Fatalf("wrong result : err is not nil. \n%s", err.Error())
	}
	device.SetupIO()
	// if err := device.StartShell(); err != nil {
	// 	t.Errorf("wrong result : err is not nil. \n%#v", err)
	// }
	if err := device.Close(); err != nil {
		t.Errorf("wrong result : err is not nil. \n%#v", err)
	}
}

func buildSSHServer(signer ssh.Signer) {
	serverConfig := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	serverConfig.AddHostKey(signer)

	target := net.JoinHostPort("0.0.0.0", testPort)
	listener, err := net.Listen("tcp", target)
	if err != nil {
		log.Fatalf("Failed to listen on %s (%s)", testPort, err)
	}
	log.Printf("Listening on %s ...", testPort)

	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept on %s (%s)", testPort, err)
		}

		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, serverConfig)
		if err != nil {
			log.Fatalf("Failed to handshake (%s)", err)
		}
		log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())

		go ssh.DiscardRequests(reqs)
		go handleChannels(chans)
	}
}

func handleChannels(chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		go handleChannel(newChannel)
	}
}

func handleChannel(newChannel ssh.NewChannel) {
	if t := newChannel.ChannelType(); t != "session" {
		if err := newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("Unknown channel type: %s", t)); err != nil {
			log.Fatal(err)
		}
		return
	}

	sshChannel, _, err := newChannel.Accept()
	if err != nil {
		log.Fatalf("Could not accept channel (%s)", err)
		return
	}

	bash := exec.Command("bash")

	close := func() {
		if e := sshChannel.Close(); e != nil {
			log.Fatal(e)
		}
		_, e := bash.Process.Wait()
		if e != nil {
			log.Printf("Failed to exit bash (%s)", e)
		}
		log.Printf("Session closed")
	}

	f, err := pty.Start(bash)
	if err != nil {
		log.Printf("Could not start pty (%s)", err)
		close()
		return
	}

	var once sync.Once
	go func() {
		if _, err := io.Copy(sshChannel, f); err != nil {
			log.Fatal(err)
		}
		once.Do(close)
	}()
	go func() {
		if _, err := io.Copy(f, sshChannel); err != nil {
			log.Fatal(err)
		}
		once.Do(close)
	}()
}
