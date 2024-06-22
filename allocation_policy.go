package lvm2go

type AllocationPolicy string

const (
	Contiguous  AllocationPolicy = "contiguous"
	Normal      AllocationPolicy = "normal"
	Cling       AllocationPolicy = "cling"
	ClingByTags AllocationPolicy = "cling_by_tags"
	Anywhere    AllocationPolicy = "anywhere"
	Inherit     AllocationPolicy = "inherit"
)

func (opt AllocationPolicy) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.AllocationPolicy = opt
}

func (opt AllocationPolicy) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AppendAll([]string{"--alloc", string(opt)})
	return nil
}
