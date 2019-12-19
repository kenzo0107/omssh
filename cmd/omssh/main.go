package main

import (
	"errors"
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

	"github.com/kenzo0107/omssh"
	"github.com/kenzo0107/omssh/pkg/awsapi"
	"github.com/kenzo0107/omssh/pkg/utility"

	latest "github.com/tcnksm/go-latest"
)

const (
	name        = "omssh"
	version     = "0.0.3"
	defaultUser = "ubuntu"
)

var (
	defUsers = []string{"ubuntu", "ec2-user"}

	flags = []cli.Flag{
		cli.StringFlag{
			Name:  "region, r",
			Value: "ap-northeast-1",
			Usage: "aws region",
		},
		cli.StringFlag{
			Name:  "port, p",
			Value: "22",
			Usage: "ssh port",
		},
		cli.BoolFlag{
			Name:  "user, u",
			Usage: "select ssh user",
		},
	}

	app = &cli.App{
		Name:    name,
		Version: version,
		Flags:   flags,
	}
)

func main() {
	app.Action = func(c *cli.Context) error {
		credentialsPath := getCredentialsPath(runtime.GOOS)
		return action(c, credentialsPath)
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getCredentialsPath(runtimeGOOS string) string {
	c := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
	if c != "" {
		return c
	}

	var configDir string
	home := os.Getenv("HOME")
	if home == "" && runtimeGOOS == "windows" {
		configDir = os.Getenv("APPDATA")
	} else {
		configDir = home
	}
	return filepath.Join(configDir, ".aws", "credentials")
}

func checkLatest(version string) error {
	version = fixVersionStr(version)
	githubTag := &latest.GithubTag{
		Owner:             "kenzo0107",
		Repository:        "omssh",
		FixVersionStrFunc: fixVersionStr,
	}
	res, err := latest.Check(githubTag, version)
	if err != nil {
		return err
	}
	if res.Outdated {
		return errors.New("not latest, you should upgrade")
	}
	return nil
}

func fixVersionStr(v string) string {
	v = strings.TrimPrefix(v, "v")
	vs := strings.Split(v, "-")
	return vs[0]
}

func action(c *cli.Context, credentialsPath string) error {
	region := c.String("region")

	profiles, err := utility.GetProfiles(credentialsPath)
	if err != nil {
		return err
	}

	profileWithAssumeRole, err := utility.FinderProfile(profiles)
	if err != nil {
		return err
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
		return err
	}

	// select an ec2
	ec2, err := awsapi.FinderEC2(ec2Instances)
	if err != nil {
		return err
	}

	user := defaultUser
	if c.Bool("user") {
		u, e := awsapi.FinderUsername(defUsers)
		if e != nil {
			return e
		}
		user = u
	}

	cache := cache.New(480*time.Minute, 1440*time.Minute)
	publicKey, privateKey := utility.SSHKeyGen(cache)

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
		return err
	}

	// ssh -i <temporary ssh private key> <user>@<public ip address>
	log.Printf("ssh %s@%s -p %s [%s]\n", user, ec2.PublicIPAddress, c.String("port"), ec2.InstanceID)

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return err
	}

	sshClientConfig := omssh.ConfigureSSHClient(user, signer)

	device := omssh.NewDevice(ec2.PublicIPAddress, c.String("port"))
	if err := device.SSHConnect(sshClientConfig); err != nil {
		return err
	}
	device.SetupIO()

	if err := device.StartShell(); err != nil {
		return err
	}
	return nil
}
