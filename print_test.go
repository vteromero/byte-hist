package main

import "testing"

func TestWrappingFormat(t *testing.T) {
	input := []struct {
		w int
		r bool
	}{
		{0, false},
		{0, true},
		{10, false},
		{10, true},
	}
	expected := []string{
		"%-0s",
		"%0s",
		"%-10s",
		"%10s",
	}

	if len(input) != len(expected) {
		t.Fatalf("the input length (%v) does not match the expected (%v)",
			len(input), len(expected))
	}

	for i, in := range input {
		out := wrappingFormat(in.w, in.r)
		if out != expected[i] {
			t.Errorf("%q does not match the expected %q", out, expected[i])
		}
	}
}

func TestHumanReadableFileSizeFmt(t *testing.T) {
	input := []struct {
		v float64
		u string
	}{
		{1.0, "bytes"},
		{68.5, "KB"},
		{712.0, "KB"},
		{18.127, "MB"},
		{44.0, "MB"},
		{999.99, "MB"},
		{999.999, "MB"},
	}
	expected := []string{
		"1 bytes",
		"68.50 KB",
		"712 KB",
		"18.13 MB",
		"44 MB",
		"999.99 MB",
		"1000.00 MB",
	}

	if len(input) != len(expected) {
		t.Fatalf("the input length (%v) does not match the expected (%v)",
			len(input), len(expected))
	}

	for i, in := range input {
		out := humanReadableFileSizeFmt(in.v, in.u)
		if out != expected[i] {
			t.Errorf("%q does not match the expected %q", out, expected[i])
		}
	}
}

func TestHumanReadableFileSize(t *testing.T) {
	input := []uint64{
		12,
		200,
		999,
		1023,
		1024,
		1555,
		2048,
		9115,
		100000,
		512000,
		2000000,
		362566123,
		16750372454,
		3846230039278812,
	}
	expected := []string{
		"12 bytes",
		"200 bytes",
		"999 bytes",
		"1023 bytes",
		"1 KB",
		"1.52 KB",
		"2 KB",
		"8.90 KB",
		"97.66 KB",
		"500 KB",
		"1.91 MB",
		"345.77 MB",
		"15.60 GB",
		"3.42 PB",
	}

	if len(input) != len(expected) {
		t.Fatalf("the input length (%v) does not match the expected (%v)",
			len(input), len(expected))
	}

	for i, in := range input {
		out := humanReadableFileSize(in)
		if out != expected[i] {
			t.Errorf("%q does not match the expected %q", out, expected[i])
		}
	}
}
