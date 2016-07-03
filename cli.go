package ghg

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/go-homedir"
	"github.com/tcnksm/go-gitconfig"
)

const (
	exitCodeOK = iota
	exitCodeParseFlagErr
	exitCodeErr
)

const version = "0.0.0"

type ghOpts struct {
	Get getCommand `description:"get stuffs" command:"get" subcommands-optional:"true"`
	Bin binCommand `description:"display bin dir" command:"bin" subcommands-optional:"true"`
}

type getCommand struct {
	targets []string
}

func (g *getCommand) Execute(args []string) error {
	bin, err := ghgBin()
	if err != nil {
		return err
	}
	ghcli := getOctCli(getToken())
	for _, target := range args {
		gh := &ghg{
			binDir: bin,
			target: target,
			client: ghcli,
		}
		err := gh.get()
		if err != nil {
			return err
		}
	}
	log.Printf("done!")
	return nil
}

func ghgHome() (string, error) {
	ghome := os.Getenv("GHG_HOME")
	if ghome != "" {
		return ghome, nil
	}
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".ghg"), nil
}

func ghgBin() (string, error) {
	home, err := ghgHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "bin"), nil
}

type binCommand struct{}

func (b *binCommand) Execute(args []string) error {
	bin, err := ghgBin()
	if err != nil {
		return err
	}
	fmt.Print(bin)
	return nil
}

// CLI is struct for command line tool
type CLI struct {
	OutStream, ErrStream io.Writer
}

// Run the ghg
func (cli *CLI) Run(argv []string) int {
	log.SetOutput(cli.ErrStream)
	log.SetFlags(0)
	err := parseArgs(argv)
	if err != nil {
		if ferr, ok := err.(*flags.Error); ok {
			if ferr.Type == flags.ErrHelp {
				return exitCodeOK
			}
			return exitCodeParseFlagErr
		}
		return exitCodeErr
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

func parseArgs(args []string) error {
	opts := &ghOpts{}
	_, err := flags.ParseArgs(opts, args)
	return err
}
