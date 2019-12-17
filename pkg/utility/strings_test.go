package utility

import "testing"

func TestConvNewline(t *testing.T) {
	if ConvNewline("a\r\nb\rc\nd", "") != "abcd" {
		t.Error("ConvNewline(\"a\\r\\nb\\rc\\nd\") should be \"abcd\", but does not match.")
	}
}

func TestConvNewlineEmptyString(t *testing.T) {
	if ConvNewline("", "") != "" {
		t.Error("ConvNewline(\"\") should be \"\", but does not match.")
	}
}
