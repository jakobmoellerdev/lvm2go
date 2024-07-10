// Package lvm2go implements a Go API for the LVM2 command line tools.
//
// The API is designed to be simple and easy to use, while still providing
// access to the full functionality of the LVM2 command line tools. (as much as possible)
//
// Compared to a simple command line wrapper, lvm2go provides a more structured
// way to interact with the LVM2 tools, and allows for more complex interactions
// with the LVM2 tools while safeguarding typing and allowing for fine-grained control
// over the input of various usually problematic parameters, such as sizes (and their conversion),
// validation of input parameters, and caching of data.
//
// A simple usage example:
//
//		func main() {
//			// Create a new LVM client
//			c, err := lvm2go.NewClient()
//			if err != nil {
//				panic(err)
//			}
//
//			// List all volume groups
//			vgs, err := c.VGs()
//			if err != nil {
//				panic(err)
//			}
//
//			// Create a new Logical Volume in the first group
//			if err = c.LVCreate(LogicalVolumeName("mylv"), VolumeGroupName(vgs[0].Name), MustParseSize("1G")); err != nil {
//				panic(err)
//			}
//	    }
package lvm2go
