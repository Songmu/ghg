package ghg

import (
	"context"
	"fmt"
)

type binCmd struct{}

func (b *binCommand) Run(ctx context.Context, args []string) error {
	bin, err := ghgBin()
	if err != nil {
		return err
	}
	fmt.Println(bin)
	return nil
}
