package lvm2go

type Thin bool

func (opt Thin) ApplyToArgs(args Arguments) error {
	if opt {
		args.AddOrReplaceAll([]string{"--thin"})
	}
	return nil
}

func (opt Thin) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Thin = opt
}
