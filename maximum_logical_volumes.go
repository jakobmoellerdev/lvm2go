package lvm2go

import (
	"fmt"
)

type MaximumLogicalVolumes int

func (opt MaximumLogicalVolumes) ApplyToVGChangeOptions(opts *VGChangeOptions) {
	opts.MaximumLogicalVolumes = opt
}

func (opt MaximumLogicalVolumes) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.MaximumLogicalVolumes = opt
}

func (opt MaximumLogicalVolumes) ApplyToArgs(args Arguments) error {
	if opt == 0 {
		return nil
	}
	switch args.GetType() {
	case ArgsTypeVGChange:
		args.AddOrReplace(fmt.Sprintf("--logicalvolume=%d", opt))
	default:
		args.AddOrReplace(fmt.Sprintf("--maxlogicalvolumes=%d", opt))
	}
	return nil
}
