package lvm2go

type PhysicalExtentSize Size

func (opt *PhysicalExtentSize) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.PhysicalExtentSize = opt
}

func (opt *PhysicalExtentSize) ApplyToArgs(args Arguments) error {
	if opt == nil {
		return nil
	}

	size := Size(*opt)

	if err := size.Validate(); err != nil {
		return err
	}

	args.AddOrReplaceAll([]string{"--physicalextentsize", size.String()})
	return nil
}
