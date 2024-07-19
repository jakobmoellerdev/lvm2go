package lvm2go

import (
	"fmt"
)

const (
	DoNotZeroVolume Zero = "n"
	ZeroVolume      Zero = "y"
)

type Zero string

func (opt Zero) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Zero = opt
}

func (opt Zero) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplace(fmt.Sprintf("--zero=%s", string(opt)))
	return nil
}
