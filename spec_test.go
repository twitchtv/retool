package main

import (
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver"
)

func TestSpecVersioning(t *testing.T) {
	readTest := func(file string, wantVersion *semver.Version) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			path := filepath.Join("testdata", "specs", file)
			spec, err := readPath(path)
			if err != nil {
				t.Fatalf("unable to read spec file: %s", err)
			}

			if wantVersion == nil {
				if spec.RetoolVersion != nil {
					t.Errorf("unexpected spec retool version, have=%q want=nil", spec.RetoolVersion, nil)
				}
			} else if !spec.RetoolVersion.Equal(wantVersion) {
				t.Errorf("unexpected spec retool version, have=%q want=%q", spec.RetoolVersion, wantVersion)
			}
		}
	}

	t.Run("read", func(t *testing.T) {
		t.Parallel()
		t.Run("unversioned", readTest("unversioned.json", nil))
		t.Run("v1.2.0", readTest("v1.2.0.json", semver.MustParse("v1.2.0")))
	})
}
