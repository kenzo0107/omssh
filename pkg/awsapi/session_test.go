package awsapi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewSession(t *testing.T) {
	fname := filepath.Join("..", "..", "testdata", "credentials")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", fname)

	for _, testcase := range []struct {
		name string
		call func(t *testing.T)
	}{
		{
			"SetProfile",
			func(t *testing.T) {
				profile := "hoge"
				region := "ap-northeast-1"
				s := newSession(profile, region)
				if e, a := region, *s.Config.Region; e != a {
					t.Errorf("expect %v, got %v", e, a)
				}

				c, err := s.Config.Credentials.Get()
				if err != nil {
					t.Error(err)
				}
				if c.AccessKeyID != "abcdefg1234567890" || c.SecretAccessKey != "abcdefghijklmnopqrstuvwxyz" {
					t.Errorf("wrong result : AcceesKeyID %s = abcdefg1234567890 or SecretAccessKey %s = abcdefghijklmnopqrstuvwxyz", c.AccessKeyID, c.SecretAccessKey)
				}
			},
		}, {
			"SetDefaultProfile",
			func(t *testing.T) {
				profile := ""
				region := "ap-southeast-1"
				s := newSession(profile, region)
				if e, a := region, *s.Config.Region; e != a {
					t.Errorf("expect %v, got %v", e, a)
				}
				c, err := s.Config.Credentials.Get()
				if err != nil {
					t.Error(err)
				}
				if c.AccessKeyID != "default1234567890" || c.SecretAccessKey != "defaultabcdefghijklmnopqrstuvwxyz" {
					t.Errorf("wrong result : AcceesKeyID %s = default1234567890 or SecretAccessKey %s = defaultabcdefghijklmnopqrstuvwxyz", c.AccessKeyID, c.SecretAccessKey)
				}
			},
		},
	} {
		t.Run(testcase.name, testcase.call)
	}
}
