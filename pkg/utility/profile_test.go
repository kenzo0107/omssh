package utility

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
	"github.com/nsf/termbox-go"
)

var (
	testProfiles = []string{
		"default",
		"hoge",
		"moge|role_arn = arn:aws:iam::1234567890:role/stsRole|source_profile = hoge|mfa_serial = arn:aws:iam::604257609175:mfa/kenzo.tanaka",
	}
)

func TestGetProfilesByTestCredentialsPath(t *testing.T) {
	cp := filepath.Join("..", "..", "testdata", "credentials")
	profiles, err := GetProfiles(cp)
	if err != nil {
		t.Error("wrong result: \nerr is not nil")
	}
	if diff := cmp.Diff(profiles, testProfiles); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}
}

func TestGetProfilesByEmptyCredentialsPath(t *testing.T) {
	profiles, err := GetProfiles("notfound_credentials")
	if err == nil {
		t.Error("wrong result: \nerr is nil")
	}
	if len(profiles) != 0 {
		t.Error("wrong result: profiles = []string{}")
	}
}

func finderProfileTesting(t *testing.T, expectedProfile string) {
	term := fuzzyfinder.UseMockedTerminal()
	term.SetSize(60, 10)

	term.SetEvents(append(
		TermboxKeys(expectedProfile),
		termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEnter})...)
	actualProfile, err := FinderProfile(testProfiles)
	if err != nil {
		t.Error("wrong result: \nerr is not nil")
	}
	if diff := cmp.Diff(expectedProfile, actualProfile); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}
}

func TestFinderProfile(t *testing.T) {
	for _, testcase := range []struct {
		name string
		call func(t *testing.T)
	}{
		{
			"profile default",
			func(t *testing.T) {
				finderProfileTesting(t, "default")
			},
		},
		{
			"profile hoge",
			func(t *testing.T) {
				finderProfileTesting(t, "hoge")
			},
		},
		{
			"profile moge",
			func(t *testing.T) {
				finderProfileTesting(t, "moge|role_arn = arn:aws:iam::1234567890:role/stsRole|source_profile = hoge|mfa_serial = arn:aws:iam::604257609175:mfa/kenzo.tanaka")
			},
		},
		{
			"profile bar not found",
			func(t *testing.T) {
				types := "bar"
				term := fuzzyfinder.UseMockedTerminal()
				term.SetSize(60, 10)

				term.SetEvents(append(
					TermboxKeys(types),
					termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEnter})...)

				profile, err := FinderProfile(testProfiles)
				if err == nil {
					t.Error("wrong result: \nerr is nil")
				}
				if diff := cmp.Diff(profile, ""); diff != "" {
					t.Error("wrong result: \nprofile is not empty")
				}
			},
		},
	} {
		t.Run(testcase.name, testcase.call)
	}
}
