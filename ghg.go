package ghg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/mholt/archiver"
	"github.com/mitchellh/ioprogress"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func getOctCli(ctx context.Context, token string) *github.Client {
	var oauthClient *http.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		oauthClient = oauth2.NewClient(ctx, ts)
	}
	return github.NewClient(oauthClient)
}

type ghg struct {
	ghgHome string
	target  string
	client  *github.Client
	upgrade bool
}

func (gh *ghg) getGhgHome() string {
	if gh.ghgHome != "" {
		return gh.ghgHome
	}
	return "."
}

func (gh *ghg) getBinDir() string {
	return filepath.Join(gh.getGhgHome(), "bin")
}

var (
	archiveReg = regexp.MustCompile(`\.(?:zip|tgz|tar\.gz)$`)
	anyExtReg  = regexp.MustCompile(`\.[a-zA-Z0-9]+$`)
	isWindows  = runtime.GOOS == "windows"
)

func (gh *ghg) get(ctx context.Context) error {
	owner, repo, tag, err := getOwnerRepoAndTag(gh.target)
	if err != nil {
		return errors.Wrap(err, "failed to resolve target")
	}
	log.Printf("fetch the GitHub release for %s\n", gh.target)

	if tag == "" {
		url := fmt.Sprintf("https://github.com/%s/%s/releases/latest", owner, repo)
		httpCli := gh.client.Client()
		req, err := gh.client.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Accept", "application/json")
		req = req.WithContext(ctx)
		resp, err := httpCli.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusOK {
			var latestResp struct {
				TagName string `json:"tag_name"`
			}
			err := json.NewDecoder(resp.Body).Decode(&latestResp)
			if err != nil {
				return err
			}
			tag = latestResp.TagName
		} else {
			rel, _, err := gh.client.Repositories.GetLatestRelease(ctx, owner, repo)
			if err != nil {
				return errors.Wrap(err, "failed to get latest tag")
			}
			tag = rel.GetTagName()
		}
		resp.Body.Close()
	}

	goarch := runtime.GOARCH
	goos := runtime.GOOS
	format := "tar.gz"
	if goos == "windows" || goos == "darwin" {
		format = "zip"
	}
	url := fmt.Sprintf(
		"https://github.com/%[1]s/%[2]s/releases/download/%[3]s/%[2]s_%[3]s_%[4]s_%[5]s.%[6]s",
		owner, repo, tag, goos, goarch, format)
	req, err := gh.client.NewRequest("HEAD", url, nil)
	if err != nil {
		return err
	}
	httpCli := gh.client.Client()
	httpCli.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	req = req.WithContext(ctx)
	resp, err := httpCli.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	var urls []string
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		urls = append(urls, url)
	} else {
		release, _, err := gh.client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
		if err != nil {
			return errors.Wrap(err, "failed to fetch a release")
		}
		for _, asset := range release.Assets {
			name := asset.GetName()
			if strings.Contains(name, goarch) && strings.Contains(name, goos) &&
				(archiveReg.MatchString(name) || !anyExtReg.MatchString(name)) {
				urls = append(urls, asset.GetBrowserDownloadURL())
			}
		}
	}
	if len(urls) < 1 {
		return fmt.Errorf("no assets available")
	}
	err = os.MkdirAll(filepath.Join(gh.getGhgHome(), "tmp"), 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create tmpdir")
	}
	log.Printf("install %s/%s version: %s", owner, repo, tag)
	for _, url := range urls {
		err := gh.install(url)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gh *ghg) install(url string) error {
	log.Printf("download %s\n", url)
	archivePath, err := gh.download(url)
	if err != nil {
		return errors.Wrap(err, "failed to download")
	}
	tmpdir := filepath.Dir(archivePath)
	defer os.RemoveAll(tmpdir)

	bin := gh.getBinDir()
	os.MkdirAll(bin, 0755)

	if !archiveReg.MatchString(url) {
		_, repo, _, _ := getOwnerRepoAndTag(gh.target)
		name := lcs(repo, filepath.Base(archivePath))
		name = strings.Trim(name, "-_")
		if name == "" {
			name = repo
		}
		if isWindows {
			name += ".exe"
		}
		return gh.place(archivePath, filepath.Join(gh.getBinDir(), name))
	}
	workDir := filepath.Join(tmpdir, "work")
	os.MkdirAll(workDir, 0755)

	log.Printf("extract %s\n", path.Base(url))
	err = extract(archivePath, workDir)
	if err != nil {
		return errors.Wrap(err, "failed to extract")
	}

	err = gh.pickupExecutable(workDir)
	if err != nil {
		return errors.Wrap(err, "failed to pickup")
	}
	return nil
}

func (gh *ghg) place(src, dest string) error {
	if exists(dest) {
		if !gh.upgrade {
			log.Printf("%s already exists. skip installing. You can use -u flag for overwrite it", dest)
			return nil
		}
		log.Printf("%s exists. overwrite it", dest)
	}
	log.Printf("install %s\n", filepath.Base(dest))
	return os.Rename(src, dest)
}

func (gh *ghg) download(url string) (fpath string, err error) {
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
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("http response not OK. code: %d, url: %s", resp.StatusCode, url)
		return
	}
	archiveBase := path.Base(url)
	tempdir, err := ioutil.TempDir(filepath.Join(gh.getGhgHome(), "tmp"), "")
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
	f, err := os.OpenFile(filepath.Join(tempdir, archiveBase), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		err = errors.Wrap(err, "failed to open file")
		return
	}
	defer f.Close()
	progressR := progbar(resp.Body, resp.ContentLength)
	_, err = io.Copy(f, progressR)
	if err != nil {
		err = errors.Wrap(err, "failed to read response")
		return
	}
	return fpath, nil
}

func progbar(r io.Reader, size int64) io.Reader {
	bar := ioprogress.DrawTextFormatBar(40)
	f := func(progress, total int64) string {
		return fmt.Sprintf(
			"%s %s",
			bar(progress, total),
			ioprogress.DrawTextFormatBytes(progress, total))
	}
	return &ioprogress.Reader{
		Reader:   r,
		Size:     size,
		DrawFunc: ioprogress.DrawTerminalf(os.Stderr, f),
	}
}

func extract(src, dest string) error {
	base := filepath.Base(src)
	if strings.HasSuffix(base, ".zip") {
		return archiver.DefaultZip.Unarchive(src, dest)
	}
	if strings.HasSuffix(base, ".tar.gz") || strings.HasSuffix(base, ".tgz") {
		return archiver.DefaultTarGz.Unarchive(src, dest)
	}
	return fmt.Errorf("failed to extract file: %s", src)
}

var targetReg = regexp.MustCompile(`^(?:([^/]+)/)?([^@]+)(?:@(.+))?$`)

func getOwnerRepoAndTag(target string) (owner, repo, tag string, err error) {
	matches := targetReg.FindStringSubmatch(target)
	if len(matches) != 4 {
		err = fmt.Errorf("failed to get owner, repo and tag")
		return
	}
	owner = matches[1]
	repo = matches[2]
	tag = matches[3]
	if owner == "" {
		owner = repo
	}
	return
}

var executableReg = func() *regexp.Regexp {
	s := `^[a-z][-_a-zA-Z0-9]+`
	if isWindows {
		s += `\.exe`
	}
	return regexp.MustCompile(s + `$`)
}()

func (gh *ghg) pickupExecutable(src string) error {
	bindir := gh.getBinDir()
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if isExecutable(info) {
			return gh.place(path, filepath.Join(bindir, info.Name()))
		}
		return nil
	})
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func isExecutable(info os.FileInfo) bool {
	if isWindows {
		return executableReg.MatchString(info.Name())
	}
	return (info.Mode()&0111) != 0 && executableReg.MatchString(info.Name())
}

func lcs(a, b string) string {
	arunes := []rune(a)
	brunes := []rune(b)
	aLen := len(arunes)
	bLen := len(brunes)
	lengths := make([][]int, aLen+1)
	for i := 0; i <= aLen; i++ {
		lengths[i] = make([]int, bLen+1)
	}
	// row 0 and column 0 are initialized to 0 already

	for i := 0; i < aLen; i++ {
		for j := 0; j < bLen; j++ {
			if arunes[i] == brunes[j] {
				lengths[i+1][j+1] = lengths[i][j] + 1
			} else if lengths[i+1][j] > lengths[i][j+1] {
				lengths[i+1][j+1] = lengths[i+1][j]
			} else {
				lengths[i+1][j+1] = lengths[i][j+1]
			}
		}
	}

	// read the substring out from the matrix
	s := make([]rune, 0, lengths[aLen][bLen])
	for x, y := aLen, bLen; x != 0 && y != 0; {
		if lengths[x][y] == lengths[x-1][y] {
			x--
		} else if lengths[x][y] == lengths[x][y-1] {
			y--
		} else {
			s = append(s, arunes[x-1])
			x--
			y--
		}
	}
	// reverse string
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return string(s)
}
