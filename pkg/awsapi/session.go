package awsapi

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewSession : return new session
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
	return session.Must(session.NewSession(&config))
}

// GetProfileWithAssumeRole : get profile with assume role and etc...
func GetProfileWithAssumeRole(profileWithAssumeRole string) (profile, roleArn, mfaSerial, sourceProfile string) {
	_p := strings.Split(profileWithAssumeRole, "|")
	profile = _p[0]

	if len(_p) > 1 {
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
	}
	return
}
