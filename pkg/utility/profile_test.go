package utility

import (
	"flag"
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
	testProfiles = []string{"hoge", "moge"}
	real         = flag.Bool("real", false, "display the actual layout to the terminal")
)

func TestGetProfilesByTestCredentialsPath(t *testing.T) {
	cp := "../../testdata/credentials"
	profiles, err := GetProfiles(cp)
	if err != nil {
		t.Error("cannot get profiles from credentilas path")
	}
	if !reflect.DeepEqual(profiles, testProfiles) {
		t.Error("GetProfiles(\"../../testdata/credentials\") should be []string{\"hoge\", \"moge\"}, but does not match.")
	}
}

func TestGetProfilesByEmptyCredentialsPath(t *testing.T) {
	profiles, err := GetProfiles("")
	actual := err.Error()
	expected := "open : no such file or directory"
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}
	if len(profiles) != 0 {
		t.Error("wrong result: profiles = []string{}")
	}
}

func TestFinderProfileHoge(t *testing.T) {
	keys := func(str string) []termbox.Event {
		s := []rune(str)
		e := make([]termbox.Event, 0, len(s))
		for _, r := range s {
			e = append(e, termbox.Event{Type: termbox.EventKey, Ch: r})
		}
		return e
	}

	term := fuzzyfinder.UseMockedTerminal()
	term.SetSize(60, 10)

	expectedProfile := "hoge"
	term.SetEvents(append(
		keys(expectedProfile),
		termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEnter})...)

	actualProfile, err := FinderProfile(testProfiles)
	if err != nil {
		t.Error("cannot get profile")
	}
	if diff := cmp.Diff(expectedProfile, actualProfile); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}

	actual := term.GetResult()

	fname := filepath.Join("..", "..", "testdata", "finderprofiles_hoge_ui.golden")
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

func TestFinderProfileMoge(t *testing.T) {
	keys := func(str string) []termbox.Event {
		s := []rune(str)
		e := make([]termbox.Event, 0, len(s))
		for _, r := range s {
			e = append(e, termbox.Event{Type: termbox.EventKey, Ch: r})
		}
		return e
	}

	term := fuzzyfinder.UseMockedTerminal()
	term.SetSize(60, 10)

	expectedProfile := "moge"
	term.SetEvents(append(
		keys(expectedProfile),
		termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEnter})...)

	actualProfile, err := FinderProfile(testProfiles)
	if err != nil {
		t.Error("cannot get profile")
	}
	if diff := cmp.Diff(expectedProfile, actualProfile); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}

	actual := term.GetResult()

	fname := filepath.Join("..", "..", "testdata", "finderprofiles_moge_ui.golden")
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
