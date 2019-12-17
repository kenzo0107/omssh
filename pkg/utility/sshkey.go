package utility

import (
	sshkey "github.com/kenzo0107/sshkeygen"
	"github.com/patrickmn/go-cache"
)

// SSHKeyGen : ssh Key gen
func SSHKeyGen(c *cache.Cache) (publicKey string, privateKey []byte) {
	pubKey, isPubKey := c.Get("publicKey")
	priKey, isPriKey := c.Get("privateKey")

	if !isPubKey || !isPriKey {
		s := sshkey.New().KeyGen()
		publicKey = s.PublicKeyStr()
		privateKey = s.PrivateKeyBytes()
		c.Set("publicKey", publicKey, cache.NoExpiration)
		c.Set("privateKey", privateKey, cache.NoExpiration)
	} else {
		publicKey = pubKey.(string)
		privateKey = priKey.([]byte)
	}

	return
}
