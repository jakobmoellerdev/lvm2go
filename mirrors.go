package lvm2go

import (
	"strconv"
)

type Mirrors int

func (opt Mirrors) ApplyToArgs(args Arguments) error {
	if opt == 0 {
		return nil
	}
	args.AddOrReplaceAll([]string{"--mirrors", strconv.Itoa(int(opt))})
	return nil
}

func (opt Mirrors) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Mirrors = opt
}
