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

package lvm2go

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ErrDeviceAlreadyClosed = errors.New("loopback device already closed")
var ErrDeviceAlreadyOpened = errors.New("loopback device was not already opened")
var ErrNoBackingFileSet = errors.New("no backing file set")
var ErrNoDeviceSet = errors.New("no device set")

const BackingFilePattern = "loopback-%s"

// LoopbackDevice is an interface that represents a loopback device created with losetup.
// It can be used to create a loopback device with a backing file, find a free loopback device,
// open (set it up) and close (detach) the loopback device.
// for more information see man losetup.
type LoopbackDevice interface {
	Open() error
	Close() error

	FindFree() error
	SetBackingFile(file string) error

	Device() string
	Size() Size
	File() string

	IsOpen() bool
	IsClosed() bool
}

// CreateLoopbackDevice creates a loopback device with the specified size that has no backing file or device.
// Example:
//
//	dev, err := CreateLoopbackDevice("4G")
//	if err != nil {
//	  panic(err)
//	}
//	fmt.Println(dev.IsOpen()) <-- false
//	fmt.Println(dev.Device()) <-- ""
//	fmt.Println(dev.File()) <-- ""
//	if err := dev.SetBackingFile(""); err != nil {
//	  panic(err)
//	}
//	fmt.Println(dev.File()) <-- "/tmp/loopback-538104538104538104538104538104538104.img"
//	if err := dev.FindFree(); err != nil {
//	  panic(err)
//	}
//	fmt.Println(dev.Device()) <-- "/dev/loop0"
//	if err := dev.Open(); err != nil {
//	  panic(err)
//	}
//	defer dev.Close()
//	fmt.Println(dev.IsOpen()) <-- true
func CreateLoopbackDevice(size Size) (LoopbackDevice, error) {
	size, err := size.ToUnit(UnitBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert size to bytes to use with truncate: %w", err)
	}
	dev := &loopbackDevice{
		size:            size,
		fileIdGenerator: newNonDeterministicID,
		commandTimeout:  60 * time.Second,
	}
	return dev, nil
}

// NewLoopbackDevice creates a loopback device with the specified size
// and returns the loopback device name
// Example:
//
//	dev, err := NewLoopbackDevice("4G")
//	if err != nil {
//	  panic(err)
//	}
//	defer dev.Close()
//	fmt.Println(dev.IsOpen()) <-- true
//	fmt.Println(dev.Device()) <-- /dev/loop0
func NewLoopbackDevice(size Size) (LoopbackDevice, error) {
	dev, err := CreateLoopbackDevice(size)
	if err != nil {
		return nil, err
	}

	if err := dev.SetBackingFile(""); err != nil {
		return nil, err
	}

	if err := dev.FindFree(); err != nil {
		return nil, err
	}

	if err := dev.Open(); err != nil {
		return nil, err
	}

	return dev, nil
}

type loopbackDevice struct {
	file            string
	device          string
	size            Size
	sectorSize      Size
	fileIdGenerator func() (string, error)
	commandTimeout  time.Duration
	opened          bool
	closed          bool
	mu              sync.RWMutex
}

func (dev *loopbackDevice) SetFileIdGenerator(generator func() (string, error)) error {
	dev.mu.Lock()
	defer dev.mu.Unlock()

	if dev.opened {
		return ErrDeviceAlreadyOpened
	}
	if dev.closed {
		return ErrDeviceAlreadyClosed
	}

	dev.fileIdGenerator = generator

	return nil
}

func (dev *loopbackDevice) SetSectorSize(size Size) error {
	dev.mu.Lock()
	defer dev.mu.Unlock()

	if dev.opened {
		return ErrDeviceAlreadyOpened
	}
	if dev.closed {
		return ErrDeviceAlreadyClosed
	}

	size, err := size.ToUnit(UnitBytes)
	if err != nil {
		return fmt.Errorf("failed to convert size to bytes to use with losetup: %w", err)
	}

	dev.sectorSize = size

	return nil
}

func (dev *loopbackDevice) String() string {
	dev.mu.RLock()
	defer dev.mu.RUnlock()
	return fmt.Sprintf("%s(%q)", dev.device, dev.file)
}

func (dev *loopbackDevice) Device() string {
	dev.mu.RLock()
	defer dev.mu.RUnlock()
	return dev.device
}

func (dev *loopbackDevice) Size() Size {
	dev.mu.RLock()
	defer dev.mu.RUnlock()
	return dev.size
}

func (dev *loopbackDevice) File() string {
	dev.mu.RLock()
	defer dev.mu.RUnlock()
	return dev.file
}

func (dev *loopbackDevice) IsClosed() bool {
	dev.mu.RLock()
	defer dev.mu.RUnlock()
	return dev.closed
}

func (dev *loopbackDevice) IsOpen() bool {
	dev.mu.RLock()
	defer dev.mu.RUnlock()
	return dev.opened
}

func (dev *loopbackDevice) Close() error {
	dev.mu.Lock()
	defer dev.mu.Unlock()

	if dev.closed || !dev.opened {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), dev.commandTimeout)
	defer cancel()

	if err := exec.CommandContext(ctx, "losetup", "-d", dev.device).Run(); err != nil {
		if isLosetupNoSuchFileOrAddressError(err) {
			dev.closed = true
			dev.opened = false
			return nil
		}
		return err
	}

	if err := os.Remove(dev.file); err != nil {
		return fmt.Errorf("failed to remove backing file %s: %w", dev.file, err)
	}

	dev.closed = true
	dev.opened = false
	return nil
}

func (dev *loopbackDevice) FindFree() error {
	dev.mu.Lock()
	defer dev.mu.Unlock()

	if dev.opened {
		return ErrDeviceAlreadyOpened
	}
	if dev.closed {
		return ErrDeviceAlreadyClosed
	}

	ctx, cancel := context.WithTimeout(context.Background(), dev.commandTimeout)
	defer cancel()

	if dev.device != "" {
		return fmt.Errorf("loopback device already has the device %s assigned", dev.device)
	}
	command := exec.CommandContext(ctx, "losetup", "-f")
	stdErr := bytes.Buffer{}
	command.Stderr = &stdErr
	loop := bytes.Buffer{}
	command.Stdout = &loop
	err := command.Run()
	if stdErr.Len() > 0 {
		err = errors.Join(err, errors.New(stdErr.String()))
	}
	if err != nil {
		return err
	}
	dev.device = strings.TrimRight(loop.String(), "\n")
	return nil
}

func (dev *loopbackDevice) SetBackingFile(file string) error {
	dev.mu.Lock()
	defer dev.mu.Unlock()

	if dev.opened {
		return ErrDeviceAlreadyOpened
	}
	if dev.closed {
		return ErrDeviceAlreadyClosed
	}

	ctx, cancel := context.WithTimeout(context.Background(), dev.commandTimeout)
	defer cancel()

	if err := dev.setFile(file); err != nil {
		return err
	}

	if _, err := os.Stat(dev.file); err == nil {
		return fmt.Errorf("backing file %s already exists", dev.file)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check for backing file existence %s: %w", dev.file, err)
	}

	args := []string{fmt.Sprintf("--size=%v", uint64(dev.size.Val)), dev.file}
	out, err := exec.CommandContext(ctx, "truncate", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to truncate backing file: %w: %s", err, string(out))
	}
	return nil
}

func (dev *loopbackDevice) setFile(file string) error {
	if dev.file != "" {
		return fmt.Errorf("loopback device already has the backing file %s assigned", dev.file)
	}

	if file == "" {
		id, err := dev.fileIdGenerator()
		if err != nil {
			return err
		}
		file = filepath.Join(os.TempDir(), fmt.Sprintf(BackingFilePattern, id))
	}

	dev.file = file
	return nil
}

func (dev *loopbackDevice) Open() error {
	dev.mu.Lock()
	defer dev.mu.Unlock()

	if dev.opened {
		return nil
	}
	if dev.closed {
		return ErrDeviceAlreadyClosed
	}

	if dev.file == "" {
		return ErrNoBackingFileSet
	}
	if dev.device == "" {
		return ErrNoDeviceSet
	}

	ctx, cancel := context.WithTimeout(context.Background(), dev.commandTimeout)
	defer cancel()

	args := []string{dev.device, dev.file}
	if dev.sectorSize.Val > 0 {
		args = append(args, fmt.Sprintf("--sector-size=%d", uint64(dev.size.Val)))
	}

	out, err := exec.CommandContext(ctx, "losetup", args...).CombinedOutput()
	if err != nil {
		return errors.Join(err, errors.New(string(out)))
	}
	dev.opened = true
	return nil
}

func isLosetupNoSuchFileOrAddressError(err error) bool {
	exitErr := &exec.ExitError{}
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		if exitErr.Stderr != nil && strings.Contains(string(exitErr.Stderr), "RequestConfirm such device or address") {
			return true
		}
	}
	return false
}

func newNonDeterministicID() (string, error) {
	nonDeterministicHash, err := newNonDeterministicHash()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(nonDeterministicHash.Sum32())), nil
}

func newNonDeterministicHash() (hash.Hash32, error) {
	hashedTestName := fnv.New32()
	randomData := make([]byte, 32)
	if _, err := rand.Read(randomData); err != nil {
		return nil, err
	}
	if _, err := hashedTestName.Write(randomData); err != nil {
		return nil, err
	}
	return hashedTestName, nil
}
