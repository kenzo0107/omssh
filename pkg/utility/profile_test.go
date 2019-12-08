package utility

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
	"github.com/nsf/termbox-go"
)

var (
	testProfiles = []string{"default", "hoge", "moge"}
)

func TestGetProfilesByTestCredentialsPath(t *testing.T) {
	cp := filepath.Join("..", "..", "testdata", "credentials")
	profiles, err := GetProfiles(cp)
	if err != nil {
		t.Error("cannot get profiles from credentilas path")
	}
	if !reflect.DeepEqual(profiles, testProfiles) {
		t.Errorf("GetProfiles(\"%s\") should be []string{\"default\", \"hoge\", \"moge\"}, but does not match.", cp)
	}
}

func TestGetProfilesByEmptyCredentialsPath(t *testing.T) {
	profiles, err := GetProfiles("")
	if e, a := "open : no such file or directory", err.Error(); e != a {
		t.Errorf("expect %v, got %v", e, a)
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
		t.Error("cannot get profile")
	}
	if diff := cmp.Diff(expectedProfile, actualProfile); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}
	actual := term.GetResult()

	g := fmt.Sprintf("finder_profiles_%s_ui.golden", expectedProfile)
	fname := filepath.Join("..", "..", "testdata", g)
	// ioutil.WriteFile(fname, []byte(actual), 0644)
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("failed to load a golden file: %s", err)
	}
	expected := string(b)
	if runtime.GOOS == "windows" {
		expected = strings.Replace(expected, "\r\n", "\n", -1)
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}
}

func TestFinderProfile(t *testing.T) {
	term := fuzzyfinder.UseMockedTerminal()
	term.SetSize(60, 10)

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
				finderProfileTesting(t, "moge")
			},
		},
	} {
		t.Run(testcase.name, testcase.call)
	}
}
