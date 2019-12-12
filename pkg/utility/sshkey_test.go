package utility

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/patrickmn/go-cache"
)

func TestSshKeyGen(t *testing.T) {
	c := cache.New(480*time.Minute, 1440*time.Minute)
	c.Delete("publicKey")
	c.Delete("privateKey")

	publicKey, privateKey := SSHKeyGen(c)
	cachedPublicKey, cachedPrivateKey := SSHKeyGen(c)

	if diff := cmp.Diff(publicKey, cachedPublicKey); diff != "" {
		t.Errorf("wrong result : \n%s", diff)
	}
	if diff := cmp.Diff(privateKey, cachedPrivateKey); diff != "" {
		t.Errorf("wrong result : \n%s", diff)
	}
}
