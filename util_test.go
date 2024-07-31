package lvm2go_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"hash"
	"hash/fnv"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

const TestExtentBytes = 1024 * 1024 // 1MiB

var TestExtentSize = MustParseSize(fmt.Sprintf("%dB", TestExtentBytes))

var sharedTestClient Client
var sharedTestClientOnce sync.Once
var sharedTestClientKey = struct{}{}

var skipRootfulTests = flag.Bool("skip-rootful-tests", false, "Name of location to greet")

func SkipOrFailTestIfNotRoot(t *testing.T) {
	if os.Geteuid() != 0 {
		if *skipRootfulTests {
			t.Skip("Skipping test because it requires root privileges to setup its environment.")
		} else {
			t.Fatalf("Failing test because it requires root privileges to setup its environment.")
		}
	}
}

func SetTestClient(ctx context.Context, client Client) context.Context {
	return context.WithValue(ctx, sharedTestClientKey, client)
}

func GetTestClient(ctx context.Context) Client {
	if client, ok := ctx.Value(sharedTestClientKey).(Client); ok {
		return client
	}
	sharedTestClientOnce.Do(func() {
		sharedTestClient = NewClient()
	})
	return sharedTestClient
}

func NewDeterministicTestID(t *testing.T) string {
	return strconv.Itoa(int(NewDeterministicTestHash(t).Sum32()))
}

func NewDeterministicTestHash(t *testing.T) hash.Hash32 {
	hashedTestName := fnv.New32()
	_, err := hashedTestName.Write([]byte(t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	return hashedTestName
}

func NewNonDeterministicTestID(t *testing.T) string {
	return strconv.Itoa(int(NewNonDeterministicTestHash(t).Sum32()))
}

func NewNonDeterministicTestHash(t *testing.T) hash.Hash32 {
	hashedTestName := fnv.New32()
	randomData := make([]byte, 32)
	if _, err := rand.Read(randomData); err != nil {
		t.Fatal(err)
	}
	if _, err := hashedTestName.Write(randomData); err != nil {
		t.Fatal(err)
	}
	return hashedTestName
}

type LoopbackDevices []LoopbackDevice

func (t LoopbackDevices) Devices() Devices {
	var devices Devices
	for _, loop := range t {
		devices = append(devices, loop.Device())
	}
	return devices
}

func (t LoopbackDevices) PhysicalVolumeNames() PhysicalVolumeNames {
	var pvs PhysicalVolumeNames
	for _, loop := range t {
		pvs = append(pvs, PhysicalVolumeName(loop.Device()))
	}
	return pvs

}

// testLoopbackCreationSync is a mutex to synchronize the creation of loopback devices in tests
// so that they don't interfere with each other by requesting the same free loopback device
var testLoopbackCreationSync = sync.Mutex{}

func MakeTestLoopbackDevice(t *testing.T, size Size) LoopbackDevice {
	t.Helper()
	ctx := context.Background()
	testLoopbackCreationSync.Lock()
	defer testLoopbackCreationSync.Unlock()

	backingFilePath := filepath.Join(t.TempDir(), fmt.Sprintf("%s.img", NewNonDeterministicTestID(t)))

	logger := slog.With("size", size, "backingFilePath", backingFilePath)

	logger.DebugContext(ctx, "creating test loopback device ...")

	size, err := size.ToUnit(UnitBytes)
	if err != nil {
		t.Fatal(err)
	}
	size.Val = RoundUp(size.Val, TestExtentBytes) + TestExtentBytes

	loop, err := CreateLoopbackDevice(size)
	if err != nil {
		t.Fatal(err)
	}
	if err := loop.FindFree(); err != nil {
		t.Fatal(err)
	}
	if err := loop.SetBackingFile(backingFilePath); err != nil {
		t.Fatal(err)
	}
	if err := loop.Open(); err != nil {
		t.Fatal(err)
	}
	logger = logger.With("loop", loop)
	logger.DebugContext(ctx, "created test loopback device successfully")
	t.Cleanup(func() {
		logger.DebugContext(ctx, "cleaning up test loopback device")
		if err := loop.Close(); err != nil {
			t.Fatal(err)
		}
	})

	return loop
}

// RoundUp rounds up n to the nearest multiple of x
func RoundUp[T int | uint | float64](n, x T) T {
	return T(math.Ceil(float64(n)/float64(x))) * x
}

// RoundDown rounds down n to the nearest multiple of x
func RoundDown[T int | uint | float64](n, x T) T {
	return T(math.Floor(float64(n)/float64(x))) * x
}

type TestVolumeGroup struct {
	Name VolumeGroupName
	t    *testing.T
}

func MakeTestVolumeGroup(t *testing.T, options ...VGCreateOption) TestVolumeGroup {
	ctx := context.Background()
	name := VolumeGroupName(NewNonDeterministicTestID(t))
	c := GetTestClient(ctx)

	if err := c.VGCreate(ctx, append(options, name, PhysicalExtentSize(TestExtentSize))...); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := c.VGRemove(ctx, name, Force(true)); err != nil {
			t.Fatal(fmt.Errorf("failed to remove volume group: %w", err))
		}
	})

	return TestVolumeGroup{
		Name: name,
		t:    t,
	}
}

type TestLogicalVolume struct {
	Options LVCreateOptionList
}

func (lv TestLogicalVolume) LogicalVolumeName() LogicalVolumeName {
	for _, opt := range lv.Options {
		switch topt := opt.(type) {
		case LogicalVolumeName:
			return topt
		}
	}
	return ""
}

func (lv TestLogicalVolume) Size() Size {
	for _, opt := range lv.Options {
		switch topt := opt.(type) {
		case Size:
			return topt
		case Extents:
			return topt.ToSize(TestExtentBytes)
		}
	}
	return Size{}
}

func (lv TestLogicalVolume) Extents() Extents {
	for _, opt := range lv.Options {
		switch topt := opt.(type) {
		case Extents:
			return topt
		}
	}
	return Extents{}
}

func (vg TestVolumeGroup) MakeTestLogicalVolume(template TestLogicalVolume) TestLogicalVolume {
	vg.t.Helper()
	ctx := context.Background()

	var logicalVolumeName LogicalVolumeName
	if lvName := template.LogicalVolumeName(); lvName == "" {
		logicalVolumeName = LogicalVolumeName(NewNonDeterministicTestID(vg.t))
		template.Options = append(template.Options, logicalVolumeName)
	} else {
		logicalVolumeName = lvName
	}

	var sizeOption LVCreateOption
	if size := template.Size(); size.Val > 0 {
		var err error
		if size, err = size.ToUnit(UnitBytes); err != nil {
			vg.t.Fatal(err)
		}
		size.Val = RoundDown(size.Val, TestExtentBytes)
		sizeOption = size
	} else if extents := template.Extents(); extents.Val > 0 {
		sizeOption = extents
	} else {
		vg.t.Logf("RequestConfirm size specified for logical volume %s, defaulting to 100M", logicalVolumeName)
		if size.Val == 0 {
			size = MustParseSize("100M")
		}
		sizeOption = size
	}
	template.Options = append(template.Options, sizeOption)

	c := GetTestClient(ctx)
	if err := c.LVCreate(ctx, vg.Name, template.Options); err != nil {
		vg.t.Fatal(err)
	}
	vg.t.Cleanup(func() {
		if err := c.LVRemove(ctx, vg.Name, logicalVolumeName); err != nil {
			if IsLVMNotFound(err) {
				vg.t.Logf("logical volume %s not found, skipping removal", logicalVolumeName)
			} else {
				vg.t.Fatal(err)
			}
		}
	})
	return TestLogicalVolume{
		Options: template.Options,
	}
}

type test struct {
	LoopDevices                  []Size
	Volumes                      []TestLogicalVolume
	AdditionalVolumeGroupOptions []VGCreateOption
}

type testInfra struct {
	loopDevices LoopbackDevices
	volumeGroup TestVolumeGroup
	lvs         []TestLogicalVolume
}

func (test test) String() string {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(test); err != nil {
		panic(err)
	}
	return buf.String()
}

func (test test) SetupDevicesAndVolumeGroup(t *testing.T) testInfra {
	t.Helper()
	var loopDevices LoopbackDevices
	for _, size := range test.LoopDevices {
		loopDevices = append(loopDevices, MakeTestLoopbackDevice(t, size))
	}
	if loopDevices == nil {
		t.Fatal("RequestConfirm loop devices defined for infra")
	}
	devices := loopDevices.Devices()

	volumeGroup := MakeTestVolumeGroup(t, append(test.AdditionalVolumeGroupOptions, PhysicalVolumesFrom(devices...))...)

	var lvs []TestLogicalVolume
	for _, lv := range test.Volumes {
		lvs = append(lvs, volumeGroup.MakeTestLogicalVolume(lv))
	}

	return testInfra{
		loopDevices: loopDevices,
		volumeGroup: volumeGroup,
		lvs:         lvs,
	}
}
