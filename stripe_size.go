package lvm2go

type StripeSize Size

func (opt StripeSize) ApplyToArgs(args Arguments) error {
	if err := Size(opt).Validate(); err != nil {
		return err
	}

	args.AddOrReplaceAll([]string{"--stripesize", Size(opt).String()})

	return nil
}

func (opt StripeSize) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.StripeSize = opt
}
