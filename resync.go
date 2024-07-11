package lvm2go

type Resync bool

func (opt Resync) ApplyToArgs(args Arguments) error {
	if opt {
		args.AddOrReplaceAll([]string{"--resync"})
	}
	return nil
}

func (opt Resync) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.Resync = opt
}
