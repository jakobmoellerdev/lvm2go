package lvm2go

type VirtualSize Size

func (opt VirtualSize) ApplyToArgs(args Arguments) error {
	if err := Size(opt).Validate(); err != nil {
		return err
	}

	args.AddOrReplaceAll([]string{"--virtualsize", Size(opt).String()})

	return nil
}

func (opt VirtualSize) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.VirtualSize = opt
}
