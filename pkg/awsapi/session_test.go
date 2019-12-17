package awsapi

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewSession(t *testing.T) {
	fname := filepath.Join("..", "..", "testdata", "credentials")
	if err := os.Setenv("AWS_SHARED_CREDENTIALS_FILE", fname); err != nil {
		t.Error("error occured in os.Setenv(\"AWS_SHARED_CREDENTIALS_FILE\")")
	}

	for _, testcase := range []struct {
		name string
		call func(t *testing.T)
	}{
		{
			"Set profile hoge",
			func(t *testing.T) {
				profile := "hoge"
				region := "ap-northeast-1"
				s := NewSession(profile, region)
				if e, a := region, *s.Config.Region; e != a {
					t.Errorf("expect %v, got %v", e, a)
				}

				c, err := s.Config.Credentials.Get()
				if err != nil {
					t.Error(err)
				}
				if diff := cmp.Diff("abcdefg1234567890", c.AccessKeyID); diff != "" {
					t.Errorf("wrong result: \n%s", diff)
				}
				if diff := cmp.Diff("abcdefghijklmnopqrstuvwxyz", c.SecretAccessKey); diff != "" {
					t.Errorf("wrong result: \n%s", diff)
				}
			},
		},
		{
			"Set profile default",
			func(t *testing.T) {
				profile := ""
				region := "ap-southeast-1"
				s := NewSession(profile, region)
				if e, a := region, *s.Config.Region; e != a {
					t.Errorf("expect %v, got %v", e, a)
				}
				c, err := s.Config.Credentials.Get()
				if err != nil {
					t.Error(err)
				}
				if diff := cmp.Diff("default1234567890", c.AccessKeyID); diff != "" {
					t.Errorf("wrong result: \n%s", diff)
				}
				if diff := cmp.Diff("defaultabcdefghijklmnopqrstuvwxyz", c.SecretAccessKey); diff != "" {
					t.Errorf("wrong result: \n%s", diff)
				}
			},
		},
		{
			"set not found profile",
			func(t *testing.T) {
				profile := "bar"
				region := "ap-northnorthnorth-1"
				s := NewSession(profile, region)
				_, err := s.Config.Credentials.Get()
				if err == nil {
					t.Error("wrong result: \n err is not nil")
				}
			},
		},
	} {
		t.Run(testcase.name, testcase.call)
	}
}

func TestGetProfileWithAssumeRole(t *testing.T) {
	for _, testcase := range []struct {
		name string
		call func(t *testing.T)
	}{
		{
			"get profile etc... from assume role",
			func(t *testing.T) {
				profileWithAssumeRole := "hoge|role_arn=arn:aws:iam::123456789012:role/stsDevMemberRole|mfa_serial=arn:aws:iam::123456789012:mfa/kenzo.tanaka|source_profile=moge"

				profile, roleArn, mfaSerial, moge := GetProfileWithAssumeRole(profileWithAssumeRole)

				if diff := cmp.Diff("hoge", profile); diff != "" {
					t.Errorf("wrong result: \n%s", diff)
				}
				if diff := cmp.Diff("arn:aws:iam::123456789012:role/stsDevMemberRole", roleArn); diff != "" {
					t.Errorf("wrong result: \n%s", diff)
				}
				if diff := cmp.Diff("arn:aws:iam::123456789012:mfa/kenzo.tanaka", mfaSerial); diff != "" {
					t.Errorf("wrong result: \n%s", diff)
				}
				if diff := cmp.Diff("moge", moge); diff != "" {
					t.Errorf("wrong result: \n%s", diff)
				}
			},
		},
	} {
		t.Run(testcase.name, testcase.call)
	}
}
