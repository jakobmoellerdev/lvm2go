package lvm2go

type Select string

func (opt Select) ApplyToLVsOptions(opts *LVsOptions) {
	opts.Select = opt
}
func (opt Select) ApplyToVGsOptions(opts *VGsOptions) {
	opts.Select = opt
}
func (opt Select) ApplyToVGRemoveOptions(opts *VGRemoveOptions) {
	opts.Select = opt
}
func (opt Select) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.Select = opt
}
