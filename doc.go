/*
 Copyright 2024 The lvm2go Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

// Package lvm2go implements a Go API for the lvm2 command line tools.
//
// The API is designed to be simple and easy to use, while still providing
// access to the full functionality of the LVM2 command line tools.
//
// Compared to a simple command line wrapper, lvm2go provides a more structured
// way to interact with lvm2, and allows for more complex interactions while safeguarding typing
// and allowing for fine-grained control over the input of various usually problematic parameters,
// such as sizes (and their conversion), validation of input parameters, and caching of data.
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
//			if err = c.LVCreate(
//			    LogicalVolumeName("mylv"),
//		    	VolumeGroupName(vgs[0].Name),
//	    		MustParseSize("1G"),
//			); err != nil {
//				panic(err)
//			}
//	    }
package lvm2go
