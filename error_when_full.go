package lvm2go

type ErrorWhenFull bool

func (opt *ErrorWhenFull) ApplyToArgs(args Arguments) error {
	if opt == nil {
		return nil
	}

	args.AddOrReplaceAll([]string{"--errorwhenfull", map[bool]string{true: "y", false: "n"}[bool(*opt)]})
	return nil
}

func (opt *ErrorWhenFull) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.ErrorWhenFull = opt
}
