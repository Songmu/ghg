package ghg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/mholt/archiver"
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

var releaseByTagURL = octokit.Hyperlink("repos/{owner}/{repo}/releases/tags/{tag}")
var archiveReg = regexp.MustCompile(`\.(?:zip|tgz|tar\.gz)$`)

func (gh *ghg) install() error {
	owner, repo, tag, err := gh.getOwnerRepoAndTag()
	if err != nil {
		return errors.Wrap(err, "failed to resolve target")
	}
	var url *url.URL
	if tag == "" {
		url, err = octokit.ReleasesLatestURL.Expand(octokit.M{"owner": owner, "repo": repo})
	} else {
		url, err = releaseByTagURL.Expand(octokit.M{"owner": owner, "repo": repo, "tag": tag})
	}
	if err != nil {
		return errors.Wrap(err, "failed to build GitHub URL")
	}
	release, r := gh.client.Releases(url).Latest()
	if r.HasError() {
		return errors.Wrap(r.Err, "failed to fetch latest release")
	}
	tag = release.TagName
	goarch := runtime.GOARCH
	goos := runtime.GOOS
	var urls []string
	for _, asset := range release.Assets {
		name := asset.Name
		if strings.Contains(name, goarch) && strings.Contains(name, goos) && archiveReg.MatchString(name) {
			urls = append(urls, fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s", owner, repo, tag, name))
		}
	}

	for _, url := range urls {
		archivePath, err := download(url)
		if err != nil {
			return errors.Wrap(err, "failed to download")
		}
		workDir := filepath.Join(filepath.Dir(archivePath), "work")
		os.MkdirAll(workDir, 0755)
		err = extract(archivePath, workDir)
		if err != nil {
			return errors.Wrap(err, "failed to extract")
		}
		err = pickupExecutable(workDir, ".")
		if err != nil {
			return errors.Wrap(err, "failed to pickup")
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
		return archiver.Unzip(src, dest)
	}
	if strings.HasSuffix(base, ".tar.gz") || strings.HasSuffix(base, ".tgz") {
		return archiver.UntarGz(src, dest)
	}
	return fmt.Errorf("failed to extract file: %s", src)
}

var targetReg = regexp.MustCompile(`^(?:([^/]+)/)?([^@]+)(?:@(.+))?$`)

func (gh *ghg) getOwnerRepoAndTag() (owner, repo, tag string, err error) {
	if matches := targetReg.FindStringSubmatch(gh.target); len(matches) == 4 {
		owner = matches[1]
		repo = matches[2]
		tag = matches[3]
		if owner == "" {
			owner = repo
		}
	}
	return
}

var executableReg = regexp.MustCompile(`^[a-z][-_a-zA-Z0-9]+(?:\.exe)?$`)

func pickupExecutable(src, dest string) error {
	defer os.RemoveAll(src)
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if name := info.Name(); (info.Mode()&0111) != 0 && executableReg.MatchString(name) {
			return os.Rename(path, filepath.Join(dest, name))
		}
		return nil
	})
}
