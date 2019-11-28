package omssh

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"

	awsapi "github.com/kenzo0107/omssh/pkg/awsapi"
	sshkey "github.com/kenzo0107/sshkeygen"
)

const (
	defUser = "ubuntu"
)

var defCredentialsPath string

func init() {
	var configDir string
	home := os.Getenv("HOME")
	if home == "" && runtime.GOOS == "windows" {
		// WindowsでHOME環境変数が定義されていない場合
		configDir = os.Getenv("APPDATA")
	} else {
		configDir = home
	}
	defCredentialsPath = filepath.Join(configDir, ".aws", "credentials")
}

func Pre(c *cli.Context) {
	region := c.String("region")

	sess, err := awsapi.AssumeRoleWithSession(region, defCredentialsPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	// get ec2 list
	ec2Client := awsapi.NewEC2(sess)
	ec2List, _ := ec2Client.GetEC2List()

	// select an ec2 instance
	ec2, err := awsapi.FinderEC2Info(ec2List)
	if err != nil {
		log.Fatal(err)
		return
	}

	user := defUser
	if c.Bool("user") {
		u, err := awsapi.FinderUsername()
		if err != nil {
			log.Fatal(err)
			return
		}
		user = u
	}

	publicKey, privateKey := sshKeyGen()

	// use ec2 instance connect to send public key
	e := awsapi.NewEC2InstanceConnect(sess)
	e.SendSSHPubKey(user, ec2.InstanceID, publicKey, ec2.AvailabilityZone)

	// ssh -i <temporary ssh private key> <user>@<public ip address>
	fmt.Printf("ssh %s@%s -p %s [%s]\n", user, ec2.PublicIPAddress, c.String("port"), ec2.InstanceID)
	doSSH(user, ec2.PublicIPAddress, c.String("port"), privateKey)
}

func sshKeyGen() (publicKey string, privateKey []byte) {
	c := cache.New(480*time.Minute, 1440*time.Minute)

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

func doSSH(user, host, port string, privateKey []byte) {
	ce := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}

	auth := []ssh.AuthMethod{}
	signer, err := ssh.ParsePrivateKey(privateKey)
	ce(err, "signer")

	auth = append(auth, ssh.PublicKeys(signer))

	// set ssh config.
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// SSH connect.
	target := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", target, sshConfig)
	ce(err, "dial")

	session, err := client.NewSession()
	ce(err, "new session")
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	term := os.Getenv("TERM")
	err = session.RequestPty(term, 25, 80, modes)
	ce(err, "request pty")

	err = session.Shell()
	ce(err, "start shell")

	err = session.Wait()
	ce(err, "return")
}
