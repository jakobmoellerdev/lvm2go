package lvm2go

type PhysicalExtentSize Size

func (opt PhysicalExtentSize) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.PhysicalExtentSize = opt
}

func (opt PhysicalExtentSize) ApplyToArgs(args Arguments) error {
	if err := Size(opt).Validate(); err != nil {
		return err
	}

	args.AddOrReplaceAll([]string{"--physicalextentsize", Size(opt).String()})
	return nil
}
