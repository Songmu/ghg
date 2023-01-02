package ghg

import (
	"context"
	"flag"
	"io"
	"log"
)

type getCmd struct {
}

func (g *getCmd) run(ctx context.Context, args []string,
	outStream, errStream io.Writer) error {
	gHome, err := ghgHome()
	if err != nil {
		return err
	}
	var upgrade bool
	fs := flag.NewFlagSet("ghg get", flag.ContinueOnError)
	fs.BoolVar(&upgrade, "u", false, "ovwerwrite the executable even if exists")

	ghcli := getOctCli(ctx, getToken())
	for _, target := range args {
		gh := &ghg{
			ghgHome: gHome,
			target:  target,
			client:  ghcli,
			upgrade: upgrade,
		}
		err := gh.get(ctx)
		if err != nil {
			return err
		}
	}
	log.Printf("done!")
	return nil
}
