# lvm2go

Package lvm2go implements a Go API for the lvm2 command line tools.

The API is designed to be simple and easy to use, while still providing
access to the full functionality of the LVM2 command line tools.

Compared to a simple command line wrapper, lvm2go provides a more structured
way to interact with lvm2, and allows for more complex interactions while safeguarding typing
and allowing for fine-grained control over the input of various usually problematic parameters,
such as sizes (and their conversion), validation of input parameters, and caching of data.

A simple usage example is shown below:

```go
package main

import (
	"context"
	"log/slog"
	"os"

	. "github.com/jakobmoellerdev/lvm2go"
)

func main() {
	if os.Geteuid() != 0 {
		panic("panicking because it requires root privileges to setup its environment.")
	}

	ctx := context.Background()
	lvm := NewClient()

	deviceSize := MustParseSize("1G")

	losetup, err := NewLoopbackDevice(deviceSize)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		if err := losetup.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	pvs := PhysicalVolumesFrom(losetup.Device())
	vgName := VolumeGroupName("test")
	lvName := LogicalVolumeName("test")
	lvSize := MustParseSize("100M")

	if err := lvm.VGCreate(ctx, vgName, pvs); err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		if err := lvm.VGRemove(ctx, vgName); err != nil {
			slog.Error(err.Error())
		}
	}()

	if err := lvm.LVCreate(ctx, vgName, lvName, lvSize); err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		if err := lvm.LVRemove(ctx, vgName, lvName); err != nil {
			slog.Error(err.Error())
		}
	}()

}
```