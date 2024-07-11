package lvm2go

type Compression bool

func (opt *Compression) ApplyToArgs(args Arguments) error {
	if opt == nil {
		return nil
	}
	args.AddOrReplaceAll([]string{"--compression", map[bool]string{true: "y", false: "n"}[bool(*opt)]})
	return nil
}

func (opt *Compression) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.Compression = opt
}
