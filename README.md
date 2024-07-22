# lvm2go (Alpha)

[![Go Reference](https://pkg.go.dev/badge/github.com/jakobmoellerdev/lvm2go.svg)](https://pkg.go.dev/github.com/jakobmoellerdev/lvm2go)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakobmoellerdev/lvm2go)](https://goreportcard.com/report/github.com/jakobmoellerdev/lvm2go)
[![License](https://img.shields.io/github/license/jakobmoellerdev/lvm2go)](https://github.com/jakobmoellerdev/lvm2go)

Package lvm2go implements a Go API for the lvm2 command line tools.

_This project is in Alpha stage and should not be used in production installations. Not all commands have been properly implemented and tested._

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
	"errors"
	"log/slog"
	"os"

	. "github.com/jakobmoellerdev/lvm2go"
)

func main() {
	if os.Geteuid() != 0 {
		panic("panicking because lvm2 requires root privileges for most operations.")
	}
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() (err error) {
	ctx := context.Background()
	lvm := NewClient()
	vgName := VolumeGroupName("test")
	lvName := LogicalVolumeName("test")
	deviceSize := MustParseSize("1G")
	lvSize := MustParseSize("100M")

	var losetup LoopbackDevice
	if losetup, err = NewLoopbackDevice(deviceSize); err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, losetup.Close())
	}()

	if err = lvm.VGCreate(ctx, vgName, PhysicalVolumesFrom(losetup.Device())); err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, lvm.VGRemove(ctx, vgName))
	}()

	if err = lvm.LVCreate(ctx, vgName, lvName, lvSize); err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, lvm.LVRemove(ctx, vgName, lvName))
	}()

	return
}
```
