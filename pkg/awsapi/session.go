package awsapi

import (
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/kenzo0107/omssh/pkg/utility"
)

// NewSession : new session specified profile
func NewSession(profile, region string) *session.Session {
	var config aws.Config
	if profile != "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{
			Region:      aws.String(region),
			Credentials: creds,
		}
	} else {
		config = aws.Config{
			Region: aws.String(region),
		}
	}
	return session.New(&config)
}

// AssumeRoleWithSession : returns switched role session from argument session and IAM
func AssumeRoleWithSession(region, defCredentialsPath string) (*session.Session, error) {
	profileWithAssumeRole, err := utility.GetProfile(defCredentialsPath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_p := strings.Split(profileWithAssumeRole, "|")
	profile := _p[0]

	var sess *session.Session

	if len(_p) > 1 {
		var roleArn string
		var mfaSerial string
		var sourceProfile string

		for _, t := range _p {
			f := strings.Split(t, "=")
			if len(f) < 2 {
				continue
			}

			k := strings.TrimSpace(f[0])
			v := strings.TrimSpace(f[1])

			switch k {
			case "role_arn":
				roleArn = v
			case "mfa_serial":
				mfaSerial = v
			case "source_profile":
				sourceProfile = v
			}
		}

		sourceSess := NewSession(sourceProfile, region)

		creds := stscreds.NewCredentials(sourceSess, roleArn, func(o *stscreds.AssumeRoleProvider) {
			o.Duration = time.Hour
			o.RoleSessionName = sourceProfile
			o.SerialNumber = aws.String(mfaSerial)
			o.TokenProvider = stscreds.StdinTokenProvider
		})

		config := aws.Config{Region: aws.String(region), Credentials: creds}

		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Config:  config,
			Profile: profile,
		}))
	} else {
		sess = NewSession(profile, region)
	}
	return sess, nil
}
