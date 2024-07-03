package lvm2go

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
)

func TestLVs(t *testing.T) {
	FailTestIfNotRoot(t)
	slog.SetDefault(slog.New(NewContextPropagatingSlogHandler(NewTestingHandler(t))))
	slog.SetLogLoggerLevel(slog.LevelDebug)

	type test struct {
		loopDevices []Size
		lvs         []Size
	}

	descriptionForTest := func(test test) string {
		t.Helper()
		totalLoopSize := 0.0
		totalLVSize := 0.0
		for _, size := range test.loopDevices {
			sizeBytes, err := size.ToUnit(UnitBytes)
			if err != nil {
				t.Error(err)
			}
			totalLoopSize += sizeBytes.Val
		}
		for _, size := range test.lvs {
			sizeBytes, err := size.ToUnit(UnitBytes)
			if err != nil {
				t.Error(err)
			}
			totalLVSize += sizeBytes.Val
		}
		loopSize, err := MustParseSize(fmt.Sprintf("%fB", totalLoopSize)).ToUnit(UnitGiB)
		if err != nil {
			t.Error(err)
		}
		lvSize, err := MustParseSize(fmt.Sprintf("%fB", totalLVSize)).ToUnit(UnitGiB)
		if err != nil {
			t.Error(err)
		}
		return fmt.Sprintf("loopCount=%v,loopSize=%v,lvCount=%v,lvSize=%v",
			len(test.loopDevices), loopSize, len(test.lvs), lvSize)
	}

	type deviceInfra struct {
		loopDevices TestLoopbackDevices
		volumeGroup TestVolumeGroup
		lvs         []TestLogicalVolume
	}

	deviceInfraForTest := func(test test) deviceInfra {
		var loopDevices TestLoopbackDevices
		for _, size := range test.loopDevices {
			loopDevices = append(loopDevices, MakeTestLoopbackDevice(t, size))
		}
		devices := loopDevices.Devices()

		volumeGroup := MakeTestVolumeGroup(t, devices...)

		var lvs []TestLogicalVolume
		for _, size := range test.lvs {
			lvs = append(lvs, volumeGroup.MakeTestLogicalVolume(size))
		}

		return deviceInfra{
			loopDevices: loopDevices,
			volumeGroup: volumeGroup,
			lvs:         lvs,
		}
	}

	for _, testCase := range []test{
		{
			loopDevices: []Size{
				MustParseSize("1G"),
			},
			lvs: []Size{
				MustParseSize("100M"),
			},
		},
		{
			loopDevices: []Size{
				MustParseSize("100M"),
			},
			lvs: []Size{
				MustParseSize("100M"),
			},
		},
		{
			loopDevices: []Size{
				MustParseSize("1G"),
			},
			lvs: []Size{
				MustParseSize("1G"),
			},
		},
		{
			loopDevices: []Size{
				MustParseSize("4G"),
			},
			lvs: []Size{
				MustParseSize("2G"),
				MustParseSize("2G"),
			},
		},
	} {
		desc := descriptionForTest(testCase)
		t.Run(desc, func(t *testing.T) {
			FailTestIfNotRoot(t)

			infra := deviceInfraForTest(testCase)

			actualLVs, err := NewClient().LVs(context.Background(), infra.volumeGroup.Name)
			if err != nil {
				t.Fatal(err)
			}

			if len(actualLVs) != len(testCase.lvs) {
				t.Fatalf("Expected %d logical volumes, got %d", len(testCase.lvs), len(actualLVs))
			}

			for _, expectedLV := range infra.lvs {
				found := false
				for _, lv := range actualLVs {
					if lv.Name != expectedLV.Name {
						continue
					}
					found = true
					break
				}
				if !found {
					t.Fatalf("Expected logical volume %s not found in LVs report", expectedLV)
				}
				if eq, err := expectedLV.Size.IsEqualTo(expectedLV.Size); err != nil || !eq {
					if err != nil {
						t.Fatalf("Size inconsistency: %s", err)
					}
				}
			}
		})
	}
}

func FailTestIfNotRoot(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Fatalf("Failing test because it requires root privileges to setup its environment.")
	}
}
