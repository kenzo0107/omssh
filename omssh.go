package omssh

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/patrickmn/go-cache"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"

	awsapi "github.com/kenzo0107/omssh/pkg/awsapi"
	"github.com/kenzo0107/omssh/pkg/utility"
	sshkey "github.com/kenzo0107/sshkeygen"
)

const (
	defUser = "ubuntu"
)

var (
	credentialsPath string
	defUsers        = []string{"ubuntu", "ec2-user"}
)

func init() {
	credentialsPath = os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
	if credentialsPath == "" {
		var configDir string
		home := os.Getenv("HOME")
		if home == "" && runtime.GOOS == "windows" {
			configDir = os.Getenv("APPDATA")
		} else {
			configDir = home
		}
		credentialsPath = filepath.Join(configDir, ".aws", "credentials")
	}
}

// Pre ... pre action of omssh
func Pre(c *cli.Context) {
	region := c.String("region")

	profiles, err := utility.GetProfiles(credentialsPath)
	if err != nil {
		return
	}

	profileWithAssumeRole, err := utility.FinderProfile(profiles)
	if err != nil {
		return
	}

	_p := strings.Split(profileWithAssumeRole, "|")

	var sess *session.Session
	if len(_p) > 1 {
		profile, roleArn, mfaSerial, sourceProfile := awsapi.GetProfileWithAssumeRole(profileWithAssumeRole)

		sourceSess := awsapi.NewSession(sourceProfile, region)

		f := func(o *stscreds.AssumeRoleProvider) {
			o.Duration = time.Hour
			o.RoleSessionName = sourceProfile
			o.SerialNumber = aws.String(mfaSerial)
			o.TokenProvider = stscreds.StdinTokenProvider
		}

		creds := stscreds.NewCredentials(sourceSess, roleArn, f)

		config := aws.Config{
			Region:      aws.String(region),
			Credentials: creds,
		}

		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Config:  config,
			Profile: profile,
		}))
	} else {
		profile := _p[0]
		sess = awsapi.NewSession(profile, region)
	}

	// get list of ec2 instances
	ec2Client := awsapi.NewEC2Client(ec2.New(sess))
	ec2Instances, err := ec2Client.DescribeRunningEC2s()
	if err != nil {
		log.Fatal(err)
		return
	}

	// select an ec2
	ec2, err := awsapi.FinderEC2(ec2Instances)
	if err != nil {
		log.Fatal(err)
		return
	}

	user := defUser
	if c.Bool("user") {
		u, e := awsapi.FinderUsername(defUsers)
		if e != nil {
			log.Fatal(e)
			return
		}
		user = u
	}

	publicKey, privateKey := sshKeyGen()

	// use ec2 instance connect to send public key
	ec2instanceconnectSvc := ec2instanceconnect.New(sess)

	input := ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: aws.String(ec2.AvailabilityZone),
		InstanceId:       aws.String(ec2.InstanceID),
		InstanceOSUser:   aws.String(user),
		SSHPublicKey:     aws.String(publicKey),
	}

	ec2InstanceConnectClient := awsapi.NewEC2InstanceConnectClient(ec2instanceconnectSvc)
	r, err := ec2InstanceConnectClient.SendSSHPubKey(input)

	if err != nil || !r {
		log.Fatal(err)
		return
	}

	// ssh -i <temporary ssh private key> <user>@<public ip address>
	log.Printf("ssh %s@%s -p %s [%s]\n", user, ec2.PublicIPAddress, c.String("port"), ec2.InstanceID)
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

	defer func() {
		err = session.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

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
