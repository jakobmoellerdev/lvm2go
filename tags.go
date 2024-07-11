package lvm2go

const TagSymbol = "@"

type Tags []string

func (opt Tags) ApplyToLVsOptions(opts *LVsOptions) {
	opts.Tags = opt
}
func (opt Tags) ApplyToVGsOptions(opts *VGsOptions) {
	opts.Tags = opt
}
func (opt Tags) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.Tags = opt
}
func (opt Tags) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Tags = opt
}
func (opt Tags) ApplyToVGRemoveOptions(opts *VGRemoveOptions) {
	opts.Tags = opt
}
func (opt Tags) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.Tags = opt
}

func (opt Tags) ApplyToArgs(args Arguments) error {
	if len(opt) == 0 {
		return nil
	}

	switch args.GetType() {
	case ArgsTypeVGCreate:
		fallthrough
	case ArgsTypeLVCreate:
		tagArgs := make([]string, 0, len(opt)*2)
		for _, tag := range opt {
			tagArgs = append(tagArgs, "--addtag", SymboledTag(tag))
		}
		args.AddOrReplaceAll(tagArgs)
	default:
		tagArgs := make([]string, 0, len(opt))
		for _, tag := range opt {
			tagArgs = append(tagArgs, SymboledTag(tag))
		}
		args.AddOrReplaceAll(tagArgs)
	}
	return nil
}

func SymboledTag(tag string) string {
	if len(tag) == 0 {
		return tag
	}
	if tag[0] != TagSymbol[0] {
		return TagSymbol + tag
	}
	return tag
}

type DelTags Tags

func (opt DelTags) ApplyToArgs(args Arguments) error {
	tagArgs := make([]string, 0, len(opt)*2)
	for _, tag := range opt {
		tagArgs = append(tagArgs, "--deltag", SymboledTag(tag))
	}
	args.AddOrReplaceAll(tagArgs)
	return nil
}

func (opt DelTags) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.DelTags = opt
}

func (opt DelTags) ApplyToVGChangeOptions(opts *VGChangeOptions) {
	opts.DelTags = opt
}
