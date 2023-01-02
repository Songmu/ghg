package ghg

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/Songmu/gitconfig"
	"github.com/jessevdk/go-flags"
)

func Run(ctx context.Context, argv []string, outStream, errStream io.Writer) error {
	log.SetOutput(errStream)
	log.SetPrefix("[ghg] ")
	fs := flag.NewFlagSet(
		fmt.Sprintf("ghg (v%s rev:%s)", version, revision), flag.ContinueOnError)

	ver := fs.Bool("version", false, "display version")
	if err := fs.Parse(argv); err != nil {
		return err
	}
	if *ver {
		return printVersion(outStream)
	}

	argv = fs.Args()
	if len(argv) < 1 {
		return errors.New("no subcommand specified")
	}
	rnr, ok := dispatch[argv[0]]
	if !ok {
		return fmt.Errorf("unknown subcommand: %s", argv[0])
	}
	return rnr.run(ctx, argv[1:], outStream, errStream)
}

var dispatch = map[string]runner{
	"bin": &binCmd{},
	"get": &getCmd{},
}

type runner interface {
	run(context.Context, []string, io.Writer, io.Writer) error
}

const (
	exitCodeOK = iota
	exitCodeParseFlagErr
	exitCodeErr
)

type ghOpts struct {
	Get getCommand `description:"get stuffs" command:"get" subcommands-optional:"true"`
	Bin binCommand `description:"display bin dir" command:"bin" subcommands-optional:"true"`
	Ver verCommand `description:"display version" command:"version" subcommands-optional:"true"`
}

type getCommand struct {
	targets []string
	Upgrade bool `short:"u" description:"overwrite the executable even if exists"`
}

func (g *getCommand) Execute(args []string) error {
	gHome, err := ghgHome()
	if err != nil {
		return err
	}
	ctx := context.TODO()
	ghcli := getOctCli(ctx, getToken())
	for _, target := range args {
		gh := &ghg{
			ghgHome: gHome,
			target:  target,
			client:  ghcli,
			upgrade: g.Upgrade,
		}
		err := gh.get(ctx)
		if err != nil {
			return err
		}
	}
	log.Printf("done!")
	return nil
}

// EnvHome is key of enviroment variable represents ghg home
const EnvHome = "GHG_HOME"

func ghgHome() (string, error) {
	ghome := os.Getenv(EnvHome)
	if ghome != "" {
		return ghome, nil
	}
	home, err := os.UserHomeDir()
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
	fmt.Println(bin)
	return nil
}

type verCommand struct{}

func (b *verCommand) Execute(args []string) error {
	fmt.Printf("ghg version: %s (rev: %s)\n", version, revision)
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
	token, _ := gitconfig.GitHubToken("")
	return token
}

func parseArgs(args []string) error {
	opts := &ghOpts{}
	_, err := flags.ParseArgs(opts, args)
	return err
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "godzil v%s (rev:%s)\n", version, revision)
	return err
}
