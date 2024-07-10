package lvm2go

type Zero bool

var zeroMapping = map[bool]string{true: "y", false: "n"}

func (opt Zero) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Zero = opt
}

func (opt Zero) ApplyToArgs(args Arguments) error {
	args.AddOrReplaceAll([]string{"--zero", zeroMapping[bool(opt)]})
	return nil
}
