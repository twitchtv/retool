package main

import (
	"path/filepath"
	"testing"
)

func TestIsLegalFile(t *testing.T) {
	testcase := func(filename string, want bool) {
		t.Run(filename, func(t *testing.T) {
			have := isLegalFile(filename)
			if have != want {
				t.Fail()
			}
		})
	}

	testcase("license.md", true)
	testcase("license.txt", true)
	testcase("LICENSE", true)
	testcase("LICENCE", true)
	testcase("LICENSE.md", true)
	testcase(filepath.Join("pkg", "LICENSE.md"), true)
	testcase("LEGAL", true)
	testcase("README", true)
	testcase("COPYING", true)
	testcase("COPYRIGHT", true)
	testcase("UNLICENSE", true)

	testcase("picture.jpeg", false)
}

func TestKeep(t *testing.T) {
	testcase := func(filename string, want bool) {
		t.Run(filename, func(t *testing.T) {
			have := keepFile(filename)
			if have != want {
				t.Fail()
			}
		})
	}

	testcase("program.go", true)
	testcase("program_test.go", false)
	testcase(filepath.Join("pkg", "program.go"), true)
	testcase(filepath.Join("pkg", "program_test.go"), false)

	testcase("assembly.s", true)
	testcase("notassembly.as", false)
	testcase("picture.gif", false)
	testcase("LICENSE.md", true)
}
