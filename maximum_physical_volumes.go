package lvm2go

import (
	"fmt"
)

type MaximumPhysicalVolumes int

func (opt MaximumPhysicalVolumes) ApplyToVGChangeOptions(opts *VGChangeOptions) {
	opts.MaximumPhysicalVolumes = opt
}

func (opt MaximumPhysicalVolumes) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.MaximumPhysicalVolumes = opt
}

func (opt MaximumPhysicalVolumes) ApplyToArgs(args Arguments) error {
	if opt == 0 {
		return nil
	}
	args.AddOrReplace(fmt.Sprintf("--maxphysicalvolumes=%d", opt))
	return nil
}
