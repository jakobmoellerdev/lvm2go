package lvm2go

import (
	"context"
)

type LVOptions struct {
	// Name of the logical volume to retrieve.
	// If empty, all logical volumes are returned.
	Name string
	// TODO: Add select options
}

// LVs returns a list of logical volumes that match the given options.
// If no logical volumes are found, nil is returned.
// It is really just a wrapper around the `lvs --reportformat json` command.
func LVs(ctx context.Context, opts LVOptions) ([]LogicalVolume, error) {
	type lvReport struct {
		Report []struct {
			LV []LogicalVolume `json:"lv"`
		} `json:"report"`
	}

	var res = new(lvReport)

	args := []string{
		"lvs",
		opts.Name,
		"-o",
		"lv_uuid,lv_name,lv_full_name,lv_path,lv_size," +
			"lv_kernel_major,lv_kernel_minor,origin,origin_size,pool_lv,lv_tags," +
			"lv_attr,vg_name,data_percent,metadata_percent,pool_lv",
		"--units",
		"b",
		"--nosuffix",
		"--reportformat",
		"json",
	}
	err := RunLVMInto(ctx, res, args...)

	if IsLVMNotFound(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if len(res.Report) == 0 {
		return nil, nil
	}

	lvs := res.Report[0].LV

	if len(lvs) == 0 {
		return nil, nil
	}

	return lvs, nil
}
