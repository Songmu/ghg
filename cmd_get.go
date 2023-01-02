package ghg

import (
	"context"
	"flag"
	"log"
)

type getCmd struct {
}

func (g *getCmd) Run(ctx context.Context, args []string) error {
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
