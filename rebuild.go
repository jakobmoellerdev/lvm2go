package lvm2go

type Rebuild bool

func (opt Rebuild) ApplyToArgs(args Arguments) error {
	if opt {
		args.AddOrReplaceAll([]string{"--rebuild"})
	}
	return nil
}

func (opt Rebuild) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.Rebuild = opt
}
