package ghg

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/octokit/go-octokit/octokit"
	"github.com/pkg/errors"
)

type ghg struct {
	binDir  string
	target  string
	version string
	client  *octokit.Client
}

func getOctCli(token string) *octokit.Client {
	var auth octokit.AuthMethod
	if token != "" {
		auth = octokit.TokenAuth{AccessToken: token}
	}
	return octokit.NewClient(auth)
}

func (gh *ghg) install() error {
	owner, repo, err := gh.getRepoAndOwner()
	if err != nil {
		return errors.Wrap(err, "failed to resolve target")
	}
	url, err := octokit.ReleasesLatestURL.Expand(octokit.M{"owner": owner, "repo": repo})
	if err != nil {
		return errors.Wrap(err, "failed to build GitHub URL")
	}
	release, r := gh.client.Releases(url).Latest()
	if r.HasError() {
		return errors.Wrap(r.Err, "failed to fetch latest release")
	}
	tag := release.TagName
	goarch := runtime.GOARCH
	goos := runtime.GOOS
	var urls []string
	for _, asset := range release.Assets {
		name := asset.Name
		if strings.Contains(name, goarch) && strings.Contains(name, goos) {
			urls = append(urls, fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s", owner, repo, tag, name))
		}
	}

	for _, url := range urls {
		archivePath, err := download(url)
		if err != nil {
			return errors.Wrap(err, "failed to download")
		}
		workDir := filepath.Join(filepath.Dir(archivePath), "work")
		fmt.Println(workDir)
		err = extract(archivePath, workDir)
		if err != nil {
			return errors.Wrap(err, "failed to extract")
		}
	}
	return nil
}

func download(url string) (fpath string, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to create request")
		return
	}
	req.Header.Set("User-Agent", fmt.Sprintf("ghg/%s", version))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.Wrap(err, "failed to create request")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "failed to read response")
		return
	}
	archiveBase := path.Base(url)
	tempdir, err := ioutil.TempDir("", "ghg-")
	if err != nil {
		err = errors.Wrap(err, "failed to create tempdir")
		return
	}
	defer func() {
		if err != nil {
			os.RemoveAll(tempdir)
		}
	}()
	fpath = filepath.Join(tempdir, archiveBase)
	f, err := os.OpenFile(filepath.Join(tempdir, archiveBase), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err = errors.Wrap(err, "failed to open file")
		return
	}
	defer f.Close()
	_, err = f.Write(body)
	if err != nil {
		err = errors.Wrap(err, "failed to read response")
		return
	}
	return fpath, nil
}

func extract(src, dest string) error {
	base := filepath.Base(src)
	if strings.HasSuffix(base, ".zip") {
		return unzip(src, dest)
	}
	return nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	var move = func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), 0755)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return errors.Wrap(err, "failed to openfile")
			}
			defer f.Close()
			_, err = io.Copy(f, rc)
			if err != nil {
				return errors.Wrap(err, "failed to copy")
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := move(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gh *ghg) getRepoAndOwner() (owner, repo string, err error) {
	arr := strings.SplitN(gh.target, "/", 2)
	if len(arr) < 1 {
		return "", "", fmt.Errorf("target invalid")
	}
	owner = arr[0]
	repo = arr[0]
	if len(arr) > 1 {
		repo = arr[1]
	}
	return
}
