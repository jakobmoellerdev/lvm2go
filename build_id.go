package lvm2go

import (
	"fmt"
	"runtime/debug"
	"sync"
)

const DefaultModuleID = "github.com/jakobmoellerdev/lvm2go"

// ModuleID returns the module ID of the library used for identification
// in the logs and other places. It defaults to using the BuildInfo
// but can be overridden by setting ModuleID to a different value.
var ModuleID = sync.OnceValue(moduleID)

// moduleID returns the module ID of the lvm2go package.
func moduleID() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	id := bi.Main.Path
	if id == "" {
		id = DefaultModuleID
	}
	if bi.Main.Version != "" {
		id = fmt.Sprintf("%s@%s", id, bi.Main.Version)
	} else {
		id = fmt.Sprintf("%s (unknown version)", id)
	}
	return id
}
