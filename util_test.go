package lvm2go

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"hash"
	"hash/fnv"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

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

// TestLoopbackDevice is a struct that holds the loopback Device and the backing file.
// It is used to create a loopback Device for testing purposes.
type TestLoopbackDevice struct {
	Device      string
	BackingFile string
}

func MakeTestLoopbackDevice(t *testing.T, size string) TestLoopbackDevice {
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
func MakeLoopbackDevice(ctx context.Context, name, size string) (string, error) {
	command := exec.Command("losetup", "-f")
	command.Stderr = os.Stderr
	loop := bytes.Buffer{}
	command.Stdout = &loop
	err := command.Run()
	if err != nil {
		return "", err
	}
	loopDev := strings.TrimRight(loop.String(), "\n")
	out, err := exec.CommandContext(ctx, "truncate", fmt.Sprintf("--size=%s", size), name).CombinedOutput()
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

type TestVolumeGroup struct {
	Name string
	t    *testing.T
}

func MakeTestVolumeGroup(t *testing.T, devices ...string) TestVolumeGroup {
	ctx := context.Background()

	name := NewNonDeterministicTestID(t)

	logger := slog.With("name", name, "devices", devices)

	err := MakeLoopbackVG(ctx, name, devices...)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		logger.DebugContext(ctx, "cleaning up volume group ...")

		if err := exec.CommandContext(ctx, "vgremove", "-f", name).Run(); err != nil {
			t.Fatal(fmt.Errorf("failed to remove volume group: %w", err))
		}

		logger.DebugContext(ctx, "cleaned up volume group")
	})

	return TestVolumeGroup{
		Name: name,
		t:    t,
	}
}

// MakeLoopbackVG creates a VG made from loopback Device by losetup
// TODO replace with proper VGCreate call.
func MakeLoopbackVG(ctx context.Context, name string, devices ...string) error {
	logger := slog.With("name", name, "devices", devices)

	args := append([]string{name}, devices...)
	logger.DebugContext(ctx, "creating volume group ...")
	out, err := exec.CommandContext(ctx, "vgcreate", args...).CombinedOutput()
	if err != nil {
		slog.ErrorContext(ctx, "failed vgcreate", "output", string(out), "error", err)
		return err
	}
	logger.DebugContext(ctx, "created volume group")
	return nil
}

type TestLogicalVolume string

func (vg TestVolumeGroup) MakeTestLogicalVolume(size string) TestLogicalVolume {
	ctx := context.Background()
	name := NewNonDeterministicTestID(vg.t)
	logger := slog.With("name", name, "size", size, "vg", vg.Name)
	args := []string{fmt.Sprintf("-L%s", size), "-n", name, vg.Name}
	logger.DebugContext(ctx, "creating logical volume ...")
	out, err := exec.CommandContext(ctx, "lvcreate", args...).CombinedOutput()
	if err != nil {
		vg.t.Fatal(fmt.Errorf("failed lvcreate: %w (stdout: %s)", err, string(out)))
	}
	logger.DebugContext(ctx, "created logical volume")
	vg.t.Cleanup(func() {
		logger.DebugContext(ctx, "cleaning up logical volume ...")
		out, err := exec.CommandContext(ctx, "lvremove", "--force", filepath.Join(vg.Name, name)).CombinedOutput()
		if err != nil {
			vg.t.Fatal(fmt.Errorf("failed lvremove: %w (stdout: %s)", err, string(out)))
		}
		logger.DebugContext(ctx, "cleaned up logical volume")
	})
	return TestLogicalVolume(name)
}
