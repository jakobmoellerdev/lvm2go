package lvm2go

import (
	"strconv"
)

type Stripes int

func (opt Stripes) ApplyToArgs(args Arguments) error {
	if opt == 0 {
		return nil
	}
	args.AddOrReplaceAll([]string{"--stripes", strconv.Itoa(int(opt))})
	return nil
}

func (opt Stripes) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Stripes = opt
}
