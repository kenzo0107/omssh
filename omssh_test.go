package omssh

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/kenzo0107/omssh/pkg/utility"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/ssh"
)

type mockClient struct {
	*ssh.Client
}

func (m *mockClient) SSHConnect(config *ssh.ClientConfig) error {
	return nil
}

func (m *mockClient) SetupIO() {

}

func (m *mockClient) StartShell() error {
	return nil
}

func TestConfigureSSHClient(t *testing.T) {
	user := "ubuntu"
	cache := cache.New(480*time.Minute, 1440*time.Minute)
	_, privateKey := utility.SSHKeyGen(cache)
	t.Logf("%#v", privateKey)
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		t.Error("wrong result : err is not nil")
	}
	sshClientConfig := ConfigureSSHClient(user, signer)

	if diff := cmp.Diff(sshClientConfig.User, user); diff != "" {
		t.Errorf("wrong result :\n%s", diff)
	}
}

// func TestSSHConnect(t *testing.T) {
// 	d := &Device{
// 		Host: "192.168.0.10",
// 		Port: "22",
// 	}
// 	d.client = &mockClient{}
// 	user := "hoge"
// 	cache := cache.New(480*time.Minute, 1440*time.Minute)
// 	_, privateKey := utility.SSHKeyGen(cache)
// 	signer, err := ssh.ParsePrivateKey(privateKey)
// 	if err != nil {
// 		t.Error("wrong result : err is not nil")
// 	}
// 	sshClientConfig := ConfigureSSHClient(user, signer)

// 	d.SSHConnect(sshClientConfig)
// }
