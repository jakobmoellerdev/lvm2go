package lvm2go

type Partial bool

func (opt Partial) ApplyToArgs(args Arguments) error {
	if opt {
		args.AddOrReplaceAll([]string{"--partial"})
	}
	return nil
}

func (opt Partial) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.Partial = opt
}
