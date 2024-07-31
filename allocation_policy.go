package lvm2go

import (
	"fmt"
)

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
func (opt AllocationPolicy) ApplyToLVChangeOptions(opts *LVCreateOptions) {
	opts.AllocationPolicy = opt
}
func (opt AllocationPolicy) ApplyToVGChangeOptions(opts *VGChangeOptions) {
	opts.AllocationPolicy = opt
}
func (opt AllocationPolicy) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.AllocationPolicy = opt
}
func (opt AllocationPolicy) ApplyToPVMoveOptions(opts *PVMoveOptions) {
	opts.AllocationPolicy = opt
}

func (opt AllocationPolicy) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplace(fmt.Sprintf("--alloc=%s", string(opt)))
	return nil
}
