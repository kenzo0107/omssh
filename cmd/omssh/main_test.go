package main

import (
	"testing"
)

func TestFixVersionStr(t *testing.T) {
	version := "0.0.2-hogehoge"
	s := fixVersionStr(version)
	if s != "0.0.2" {
		t.Error("fixVersionStr(\"0.0.2-hogehoge\") should be \"0.0.2\", but does not match.")
	}
}
