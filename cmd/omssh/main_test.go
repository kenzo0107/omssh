package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetCredentialsPathOnWindowsWithoutHome(t *testing.T) {
	runtimeGOOS := "windows"
	if err := os.Setenv("AWS_SHARED_CREDENTIALS_FILE", ""); err != nil {
		t.Error(err)
	}
	if err := os.Setenv("HOME", ""); err != nil {
		t.Error(err)
	}
	getCredentialsPath(runtimeGOOS)
}

func TestGetCredentialsPath(t *testing.T) {
	fname := filepath.Join("..", "..", "testdata", "credentials")
	if err := os.Setenv("AWS_SHARED_CREDENTIALS_FILE", fname); err != nil {
		t.Error("error occured in os.Setenv(\"AWS_SHARED_CREDENTIALS_FILE\")")
	}
	getCredentialsPath("")
}

func TestCheckLatest(t *testing.T) {
	version := "0.0.2-hogehoge"
	if err := checkLatest(version); err != nil {
		t.Error(err)
	}
}

func TestCheckNotLatest(t *testing.T) {
	version := "0.0.0-hogehoge"
	err := checkLatest(version)
	if diff := cmp.Diff(err.Error(), "not latest, you should upgrade"); diff != "" {
		t.Errorf("wrong result : err message: %s", err.Error())
	}
}

func TestCheckLatestFailed(t *testing.T) {
	version := "moge"
	err := checkLatest(version)
	if err == nil {
		t.Error(err)
	}
}

func TestFixVersionStr(t *testing.T) {
	version := "0.0.2-hogehoge"
	s := fixVersionStr(version)
	if s != "0.0.2" {
		t.Error("fixVersionStr(\"0.0.2-hogehoge\") should be \"0.0.2\", but does not match.")
	}
}
