package ghg

import (
	"context"
	"fmt"
	"io"
)

type binCmd struct{}

func (b *binCmd) run(ctx context.Context, args []string,
	outStream, errStream io.Writer) error {
	bin, err := ghgBin()
	if err != nil {
		return err
	}
	fmt.Println(bin)
	return nil
}
