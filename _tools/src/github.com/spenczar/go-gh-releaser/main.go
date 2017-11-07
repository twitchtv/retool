package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func main() {
	var (
		token   = flag.String("token", "", "github API token to use")
		release = flag.String("release", "", "release tag to upload for")
		rawRepo = flag.String("repo", "", "repo to upload for, in org/repo form like 'spenczar/shipit'")
		force   = flag.Bool("force", false, "overwrite assets if already present")
	)
	flag.Parse()

	owner, repo, err := splitRepo(*rawRepo)
	mustNotErr(err)

	dir, err := runGox()
	mustNotErr(err)

	defer func() {
		_ = os.RemoveAll(dir)
	}()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	log.Printf("looking up release ID for %q", *release)
	id, err := releaseID(client, owner, repo, *release)
	log.Printf("release ID is %v", id)
	mustNotErr(err)

	files, err := ioutil.ReadDir(dir)
	mustNotErr(err)
	for i, bin := range files {
		log.Printf("publishing %v (%d / %d)", bin.Name(), i+1, len(files))

		f, err := os.Open(filepath.Join(dir, bin.Name()))
		mustNotErr(err)

		err = publishBin(client, owner, repo, id, f)
		if isAlreadyExistsErr(err) && *force {
			log.Printf("overwriting existing bin for %v", bin.Name())
			err = deleteBin(client, owner, repo, id, filepath.Base(f.Name()))
			mustNotErr(err)

			f, err = os.Open(filepath.Join(dir, bin.Name()))
			mustNotErr(err)
			err = publishBin(client, owner, repo, id, f)
		}
		mustNotErr(err)
	}

}

func runGox() (string, error) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}
	cmd := exec.Command("gox", "-output", filepath.Join(dir, "{{.Dir}}_{{.OS}}_{{.Arch}}"), ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return dir, err
}

func mustNotErr(err error) {
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

func releaseID(client *github.Client, owner, repo, tag string) (int, error) {
	release, _, err := client.Repositories.GetReleaseByTag(context.TODO(), owner, repo, tag)
	if err != nil {
		return 0, err
	}
	return *release.ID, nil
}

func publishBin(client *github.Client, owner, repo string, releaseID int, bin *os.File) error {
	opt := &github.UploadOptions{Name: filepath.Base(bin.Name())}
	_, _, err := client.Repositories.UploadReleaseAsset(context.TODO(), owner, repo, releaseID, opt, bin)
	return err
}

func deleteBin(client *github.Client, owner, repo string, releaseID int, name string) error {
	assets, err := listAssets(client, owner, repo, releaseID)
	if err != nil {
		return err
	}

	var assetID int = -1
	for _, a := range assets {
		if *a.Name == name {
			assetID = *a.ID
			break
		}
	}
	if assetID == -1 {
		return errors.Errorf("asset %q not found", name)
	}

	_, err = client.Repositories.DeleteReleaseAsset(context.TODO(), owner, repo, assetID)
	return err
}

func listAssets(client *github.Client, owner, repo string, releaseID int) ([]*github.ReleaseAsset, error) {
	opt := &github.ListOptions{}
	var allAssets []*github.ReleaseAsset
	for {
		assets, resp, err := client.Repositories.ListReleaseAssets(context.TODO(), owner, repo, releaseID, opt)
		if err != nil {
			return nil, err
		}
		allAssets = append(allAssets, assets...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allAssets, nil
}

func splitRepo(r string) (owner, repo string, err error) {
	split := strings.SplitN(r, "/", 2)
	if len(split) != 2 {
		return "", "", errors.Errorf("-repo should have one slash in it, have %q", r)
	}
	return split[0], split[1], nil
}

func isAlreadyExistsErr(err error) bool {
	ghErr, ok := err.(*github.ErrorResponse)
	if !ok {
		return false
	}

	if len(ghErr.Errors) != 1 {
		return false
	}

	return ghErr.Errors[0].Code == "already_exists"
}
