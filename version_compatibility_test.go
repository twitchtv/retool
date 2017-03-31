package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func assertRetoolDoWorks(t *testing.T, wd string, retoolBin string) {
	output := runRetoolCmd(t, wd, retoolBin, "do", "retool_test_app")
	want := "v1"
	if strings.TrimSpace(output) != want {
		t.Errorf("'retool do' gave unexpected output, have=%q  want=%q", output, want)
	}
}

// TestCompatibility tests that the current version of retool is compatible with
// old versions. It should set up projects that are readable by old versions,
// and it should be able to read projects set up by old versions.
func TestCompatibility(t *testing.T) {
	newRetool, err := buildRetool()
	if err != nil {
		t.Fatalf("unable to build current retool version: %s", err)
	}

	compatibilityTest := func(oldVersion string) {
		t.Run(oldVersion, func(t *testing.T) {
			t.Parallel()
			// Install old version
			installDir, installCleanup := setupTempDir(t)
			defer installCleanup()

			runRetoolCmd(t, installDir, newRetool, "add", "github.com/twitchtv/retool", oldVersion)
			oldRetool := filepath.Join(installDir, "_tools", "bin", "retool")

			t.Run("parallel tests", func(t *testing.T) {
				// The extra t.Run("parallel tests" grouping here is necessary to make

				// sure that 'defer installCleanup()' gets run *after* these two
				// subtests run.
				t.Run("project set up with "+oldVersion, func(t *testing.T) {
					t.Parallel()
					project, cleanup := setupTempDir(t)
					defer cleanup()
					runRetoolCmd(t, project, oldRetool, "add", "github.com/spenczar/retool_test_app", "v1")
					runRetoolCmd(t, project, oldRetool, "sync")

					assertRetoolDoWorks(t, project, newRetool)
				})

				t.Run("project set up with new version", func(t *testing.T) {
					t.Parallel()
					project, cleanup := setupTempDir(t)
					defer cleanup()
					runRetoolCmd(t, project, newRetool, "add", "github.com/spenczar/retool_test_app", "v1")
					runRetoolCmd(t, project, newRetool, "sync")

					assertRetoolDoWorks(t, project, oldRetool)
				})
			})
		})
	}

	compatibilityTest("v1.2.0")
	compatibilityTest("v1.1.0")
	compatibilityTest("v1.0.3")
}
