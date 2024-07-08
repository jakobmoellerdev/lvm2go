package lvm2go

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"hash"
	"hash/fnv"
	"log/slog"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
)

const TestExtentBytes = 16 * 1024 * 1024 // 16MiB

var TestExtentSize = MustParseSize(fmt.Sprintf("%dB", TestExtentBytes))

var sharedTestClient Client
var sharedTestClientOnce sync.Once
var sharedTestClientKey = struct{}{}

func FailTestIfNotRoot(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Fatalf("Failing test because it requires root privileges to setup its environment.")
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

type TestLoopbackDevices []TestLoopbackDevice

func (t TestLoopbackDevices) Devices() Devices {
	var devices Devices
	for _, loop := range t {
		devices = append(devices, loop.Device)
	}
	return devices
}

// TestLoopbackDevice is a struct that holds the loopback Device and the backing file.
// It is used to create a loopback Device for testing purposes.
type TestLoopbackDevice struct {
	Device      string
	BackingFile string
}

func MakeTestLoopbackDevice(t *testing.T, size Size) TestLoopbackDevice {
	ctx := context.Background()

	backingFilePath := filepath.Join(t.TempDir(), fmt.Sprintf("%s.img", NewNonDeterministicTestID(t)))

	logger := slog.With("size", size, "backingFilePath", backingFilePath)

	logger.DebugContext(ctx, "creating test loopback device ...")
	loop, err := MakeLoopbackDevice(ctx, backingFilePath, size)
	if err != nil {
		t.Fatal(err)
	}
	logger = logger.With("loop", loop)
	logger.DebugContext(ctx, "created test loopback device successfully")

	testDevice := TestLoopbackDevice{
		Device:      loop,
		BackingFile: backingFilePath,
	}

	t.Cleanup(func() {
		logger.DebugContext(ctx, "cleaning up test loopback device")

		if err := exec.CommandContext(ctx, "losetup", "-d", testDevice.Device).Run(); err != nil {
			t.Fatal(fmt.Errorf("failed to detach test loopback Device: %w", err))
		}
		if err := os.Remove(testDevice.BackingFile); err != nil {
			t.Fatal(fmt.Errorf("failed to remove test backing file: %w", err))
		}
	})

	return testDevice
}

// MakeLoopbackDevice creates a loopback Device with the specified size
// and returns the loopback Device name
// Example:
//
//	MakeLoopbackDevice(ctx, "/tmp/loopback.img", "4G")
//	// returns /dev/loop0
func MakeLoopbackDevice(ctx context.Context, name string, size Size) (string, error) {
	command := exec.Command("losetup", "-f")
	command.Stderr = os.Stderr
	loop := bytes.Buffer{}
	command.Stdout = &loop
	err := command.Run()
	if err != nil {
		return "", err
	}
	loopDev := strings.TrimRight(loop.String(), "\n")

	// add an extra extent to the size to account for metadata
	size, err = size.ToUnit(UnitBytes)
	if err != nil {
		return "", err
	}
	size.Val = RoundUp(size.Val, TestExtentBytes) + TestExtentBytes

	out, err := exec.CommandContext(ctx, "truncate", fmt.Sprintf("--size=%v", uint64(size.Val)), name).CombinedOutput()
	if err != nil {
		slog.ErrorContext(ctx, "failed truncate", "output", string(out), "error", err)
		return "", err
	}
	out, err = exec.CommandContext(ctx, "losetup", loopDev, name).CombinedOutput()
	if err != nil {
		slog.ErrorContext(ctx, "failed losetup", "output", string(out), "error", err)
		return "", err
	}
	return loopDev, nil
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

func MakeTestVolumeGroup(t *testing.T, devices ...string) TestVolumeGroup {
	ctx := context.Background()
	name := VolumeGroupName(NewNonDeterministicTestID(t))
	c := GetTestClient(ctx)

	if err := c.VGCreate(ctx, name, PhysicalVolumesFrom(devices...), PhysicalExtentSize(TestExtentSize)); err != nil {
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
		switch opt.(type) {
		case LogicalVolumeName:
			return opt.(LogicalVolumeName)
		}
	}
	return ""
}

func (lv TestLogicalVolume) Size() Size {
	for _, opt := range lv.Options {
		switch opt.(type) {
		case Size:
			return opt.(Size)
		}
	}
	return Size{}
}

func (vg TestVolumeGroup) MakeTestLogicalVolume(template TestLogicalVolume) TestLogicalVolume {
	ctx := context.Background()

	var logicalVolumeName LogicalVolumeName
	if lvName := template.LogicalVolumeName(); lvName == "" {
		logicalVolumeName = LogicalVolumeName(NewNonDeterministicTestID(vg.t))
		template.Options = append(template.Options, logicalVolumeName)
	} else {
		logicalVolumeName = lvName
	}

	size := template.Size()
	if size.Val == 0 {
		size = MustParseSize("1G")
	}
	var err error
	if size, err = size.ToUnit(UnitBytes); err != nil {
		vg.t.Fatal(err)
	}
	size.Val = RoundUp(size.Val, TestExtentBytes)

	c := GetTestClient(ctx)
	if err := c.LVCreate(ctx, vg.Name, logicalVolumeName, size); err != nil {
		vg.t.Fatal(err)
	}
	vg.t.Cleanup(func() {
		if err := c.LVRemove(ctx, vg.Name, logicalVolumeName); err != nil {
			vg.t.Fatal(err)
		}
	})
	return TestLogicalVolume{
		Options: template.Options,
	}
}

type test struct {
	loopDevices []Size
	lvs         []TestLogicalVolume
}

type testInfra struct {
	loopDevices TestLoopbackDevices
	volumeGroup TestVolumeGroup
	lvs         []TestLogicalVolume
}

func (test test) String() string {
	totalLoopSize := 0.0
	totalLVSize := 0.0
	for _, size := range test.loopDevices {
		sizeBytes := size.ToUnitIfValid(UnitBytes)
		totalLoopSize += sizeBytes.Val
	}
	for _, lv := range test.lvs {
		sizeBytes := lv.Size().ToUnitIfValid(UnitBytes)
		totalLVSize += sizeBytes.Val
	}
	loopSize := MustParseSize(fmt.Sprintf("%fB", totalLoopSize)).ToUnitIfValid(UnitGiB)
	lvSize := MustParseSize(fmt.Sprintf("%fB", totalLVSize)).ToUnitIfValid(UnitGiB)
	return fmt.Sprintf("loopCount=%v,loopSize=%v,lvCount=%v,lvSize=%v",
		len(test.loopDevices), loopSize, len(test.lvs), lvSize)
}

func (test test) SetupDevicesAndVolumeGroup(t *testing.T) testInfra {
	t.Helper()
	var loopDevices TestLoopbackDevices
	for _, size := range test.loopDevices {
		loopDevices = append(loopDevices, MakeTestLoopbackDevice(t, size))
	}
	if loopDevices == nil {
		t.Fatal("No loop devices defined for infra")
	}
	devices := loopDevices.Devices()

	volumeGroup := MakeTestVolumeGroup(t, devices...)

	var lvs []TestLogicalVolume
	for _, lv := range test.lvs {
		lvs = append(lvs, volumeGroup.MakeTestLogicalVolume(lv))
	}

	return testInfra{
		loopDevices: loopDevices,
		volumeGroup: volumeGroup,
		lvs:         lvs,
	}
}
