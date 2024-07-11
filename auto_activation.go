package lvm2go

type AutoActivation bool

func (opt AutoActivation) ApplyToArgs(args Arguments) error {
	args.AddOrReplaceAll([]string{"--setautoactivation", map[bool]string{true: "y", false: "n"}[bool(opt)]})
	return nil
}

func (opt AutoActivation) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.AutoActivation = opt
}
