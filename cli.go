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
	Get getCommand `description:"get stuffs" command:"get" subcommands-optional:"true"`
}

type getCommand struct {
	BinDir  string
	targets []string
}

func (g *getCommand) Execute(args []string) error {
	ghcli := getOctCli(getToken())
	for _, target := range args {
		gh := &ghg{
			binDir: g.BinDir,
			target: target,
			client: ghcli,
		}
		err := gh.install()
		if err != nil {
			return err
		}
	}
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
			log.Println(ferr.Error())
			return exitCodeParseFlagErr
		}
		log.Println(err.Error())
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
