package lvm2go

type Deduplication bool

func (opt *Deduplication) ApplyToArgs(args Arguments) error {
	if opt == nil {
		return nil
	}

	args.AddOrReplaceAll([]string{"--deduplication", map[bool]string{true: "y", false: "n"}[bool(opt)]})
	return nil
}

func (opt *Deduplication) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.Deduplication = opt
}
