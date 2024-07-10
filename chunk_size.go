package lvm2go

type ChunkSize Size

func (opt ChunkSize) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.ChunkSize = opt
}

func (opt ChunkSize) ApplyToArgs(args Arguments) error {
	if err := Size(opt).Validate(); err != nil {
		return err
	}

	args.AddOrReplaceAll([]string{"--chunksize", Size(opt).String()})
	return nil
}
