package ghg

import (
	"fmt"
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
		fmt.Println(archivePath)
	}
	return nil
}

func download(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("User-Agent", fmt.Sprintf("ghg/%s", version))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to create request")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response")
	}
	archiveBase := path.Base(url)
	tempdir, err := ioutil.TempDir("", "ghg-")
	if err != nil {
		return "", errors.Wrap(err, "failed to create tempdir")
	}
	fpath := filepath.Join(tempdir, archiveBase)
	f, err := os.OpenFile(filepath.Join(tempdir, archiveBase), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", errors.Wrap(err, "failed to open file")
	}
	defer f.Close()
	_, err = f.Write(body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response")
	}
	return fpath, nil
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
