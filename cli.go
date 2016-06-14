package ghg

import (
	"io"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/tcnksm/go-gitconfig"
)

const (
	exitCodeOK = iota
	exitCodeParseFlagErr
	exitCodeErr
)

const version = "0.0.0"

type ghOpts struct {
	BinDir  string
	targets []string
}

// CLI is struct for command line tool
type CLI struct {
	OutStream, ErrStream io.Writer
}

// Run the ghg
func (cli *CLI) Run(argv []string) int {
	log.SetOutput(cli.ErrStream)
	p, opts, err := parseArgs(argv)
	if err != nil {
		if ferr, ok := err.(*flags.Error); !ok || ferr.Type != flags.ErrHelp {
			p.WriteHelp(cli.ErrStream)
		}
		return exitCodeParseFlagErr
	}

	ghcli := getOctCli(getToken())
	for _, target := range opts.targets {
		gh := &ghg{
			binDir: opts.BinDir,
			target: target,
			client: ghcli,
		}
		err := gh.install()
		if err != nil {
			log.Println(err.Error())
		}
	}
	return exitCodeOK
}

func getToken() string {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		return token
	}
	token, _ = gitconfig.GithubToken()
	return token
}

func parseArgs(args []string) (*flags.Parser, *ghOpts, error) {
	opts := &ghOpts{}
	p := flags.NewParser(opts, flags.Default)
	p.Usage = "[OPTIONS]\n\nVersion: " + version
	rest, err := p.ParseArgs(args)
	opts.targets = rest
	return p, opts, err
}
