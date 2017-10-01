package goarmorlogs

import (
	"regexp"
	"testing"
)

func TestDirectoryAndLogName(t *testing.T) {
	var a = []struct {
		in, dir, file string
	}{
		{"/var/log/app.log", "/var/log", "app.log"},
		{"/var/log", "/var/log", ""},
		{"/var/log/", "/var/log", ""},
		{"", ".", ""},
		{"/", "/", ""},
		{"//", "/", ""},
		{"///", "/", ""},
		{"///...", "/", ""},
		{"..///", "..", ""},
		{"...///", "...", ""},
		{"....///", "....", ""},
		{"...", ".", ""},
		{".", ".", ""},
	}

	var (
		x, y string
		r    *regexp.Regexp

		err error
	)

	r, err = regexp.Compile("^[.]+$")
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range a {
		x, y, err = DirectoryAndLogName(v.in, r)
		if err != nil || x != v.dir || y != v.file {
			t.Errorf("directoryAndLogName(%q) => (%q, %q, %v) want (%q, %q, <nil>)",
				v.in, x, y, err, v.dir, v.file)
		}
	}
}

func TestPartitionedPathByUserID(t *testing.T) {
	var a = []struct {
		in   int64
		want string
	}{
		{1, "000/000/001"},
		{123456789, "123/456/789"},
	}

	var (
		s   string
		err error
	)

	for _, v := range a {
		s, err = PartitionedPathByUserID(v.in)
		if err != nil || s != v.want {
			t.Errorf("PartitionedPathByUserID(%d) => (%q, %v) want (%q, <nil>)",
				v.in, s, err, v.want)
		}
	}
}
