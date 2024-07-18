package lvm2go

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"time"
)

type Version struct {
	LVMVersion     string
	LibraryVersion string
	DriverVersion  string

	LVMBuild     time.Time
	LibraryBuild time.Time

	ConfigurationFlags []string
}

type (
	VersionOptions struct{}
	VersionOption  interface {
		ApplyToVersionOptions(opts *VersionOptions)
	}
	VersionOptionsList []VersionOption
)

var (
	_ ArgumentGenerator = VersionOptionsList{}
	_ Argument          = (*VersionOptions)(nil)
	_ VersionOption     = (*VersionOptions)(nil)
)

func DefaultVersionOutputProcessor() (*Version, RawOutputProcessor) {
	version := Version{}
	return &version, func(line io.Reader) error {
		scanner := bufio.NewScanner(line)
		versionLine := scanner.Scan()
		if !versionLine {
			return fmt.Errorf("no version line found")
		}
		versionFields := strings.Fields(strings.TrimSpace(scanner.Text()))
		if len(versionFields) < 4 {
			return fmt.Errorf("invalid version line: %q", scanner.Text())
		}
		version.LVMVersion = versionFields[2]
		if lvmBuildDate, err := time.Parse(time.DateOnly, strings.Trim(versionFields[3], "()")); err != nil {
			return fmt.Errorf("failed to parse library build date: %v", err)
		} else {
			version.LVMBuild = lvmBuildDate
		}

		libraryLine := scanner.Scan()
		if !libraryLine {
			return fmt.Errorf("no library line found")
		}
		libraryFields := strings.Fields(strings.TrimSpace(scanner.Text()))
		if len(libraryFields) < 4 {
			return fmt.Errorf("invalid version line: %q", scanner.Text())
		}
		version.LibraryVersion = libraryFields[2]

		if libraryBuildDate, err := time.Parse(time.DateOnly, strings.Trim(libraryFields[3], "()")); err != nil {
			return fmt.Errorf("failed to parse library build date: %v", err)
		} else {
			version.LibraryBuild = libraryBuildDate
		}

		driverLine := scanner.Scan()
		if !driverLine {
			return fmt.Errorf("no driver line found")
		}
		driverFields := strings.Fields(strings.TrimSpace(scanner.Text()))
		if len(driverFields) < 3 {
			return fmt.Errorf("invalid version line: %q", scanner.Text())
		}
		version.DriverVersion = driverFields[2]

		configurationLine := scanner.Scan()
		if !configurationLine {
			return fmt.Errorf("no configuration line found")
		}
		configurationFields := strings.Fields(strings.TrimSpace(scanner.Text()))
		if len(configurationFields) < 3 {
			return fmt.Errorf("invalid version line: %q", scanner.Text())
		}
		version.ConfigurationFlags = configurationFields[2:]
		return scanner.Err()
	}
}

func (c *client) Version(ctx context.Context, opts ...VersionOption) (Version, error) {
	args, err := VersionOptionsList(opts).AsArgs()
	if err != nil {
		return Version{}, err
	}

	version, versionProcessor := DefaultVersionOutputProcessor()

	if err := c.RunLVMRaw(ctx, versionProcessor, append([]string{"version"}, args.GetRaw()...)...); err != nil {
		return Version{}, fmt.Errorf("failed to get version: %v", err)
	}

	return *version, nil
}

func (list VersionOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := VersionOptions{}
	for _, opt := range list {
		opt.ApplyToVersionOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *VersionOptions) ApplyToVersionOptions(new *VersionOptions) {
	*new = *opts
}

func (opts *VersionOptions) ApplyToArgs(Arguments) error {
	return nil
}
