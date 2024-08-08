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

package lvm2go_test

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
	. "github.com/jakobmoellerdev/lvm2go/config"
)

func Test_RawConfig(t *testing.T) {
	t.Parallel()
	SkipOrFailTestIfNotRoot(t)
	slog.SetDefault(slog.New(NewContextPropagatingSlogHandler(NewTestingHandler(t))))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx := context.Background()
	clnt := GetTestClient(ctx)

	ver, err := clnt.RawConfig(ctx, ConfigTypeFull)

	if err != nil {
		t.Fatalf("failed to get config: %v", err)
	}

	if len(ver) == 0 {
		t.Fatalf("RawConfig is empty")
	}

	profileDir, err := GetFromRawConfig[string](ver, "profile_dir")
	if err != nil {
		t.Fatalf("failed to get profile_dir: %v", err)
	}

	if len(profileDir) == 0 {
		t.Fatalf("profile dir is empty even though that was not expected")
	}
}

func Test_DecodeConfig(t *testing.T) {
	t.Parallel()
	SkipOrFailTestIfNotRoot(t)
	slog.SetDefault(slog.New(NewContextPropagatingSlogHandler(NewTestingHandler(t))))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx := context.Background()
	clnt := GetTestClient(ctx)

	t.Run("all structs", func(t *testing.T) {
		type structConfig struct {
			Config struct {
				ProfileDir string `lvm:"profile_dir"`
			} `lvm:"config"`
		}
		c := &structConfig{}
		if err := clnt.ReadAndDecodeConfig(ctx, c, ConfigTypeFull); err != nil {
			t.Fatalf("failed to get config: %v", err)
		}
		if len(c.Config.ProfileDir) == 0 {
			t.Fatalf("profile dir is empty even though that was not expected")
		}
	})

	t.Run("all pointers", func(t *testing.T) {
		type pointerConfig struct {
			Config *struct {
				ProfileDir *string `lvm:"profile_dir"`
			} `lvm:"config"`
			NoPoint *struct{} `lvm:"no-point"`
		}
		c := &pointerConfig{}
		if err := clnt.ReadAndDecodeConfig(ctx, c, ConfigTypeFull); err != nil {
			t.Fatalf("failed to get config: %v", err)
		}
		if c.Config == nil || c.Config.ProfileDir == nil || len(*c.Config.ProfileDir) == 0 {
			t.Fatalf("profile dir is empty even though that was not expected")
		}
	})

	t.Run("bad inner block", func(t *testing.T) {
		type ignoredOuterBlockConfig struct {
			Config struct {
				ProfileDir string `lvm:"profile_dir"`
			} `lvm:"config"`
			Bla string
		}
		c := &ignoredOuterBlockConfig{}
		if err := clnt.ReadAndDecodeConfig(ctx, c, ConfigTypeFull); err == nil {
			t.Fatalf("expected error due to config block requiring struct or pointer to struct")
		}
	})

	t.Run("bad inner block ignored", func(t *testing.T) {
		type ignoredOuterBlockConfig struct {
			Config struct {
				ProfileDir string `lvm:"profile_dir"`
			} `lvm:"config"`
			Bla string `lvm:"-"`
		}
		c := &ignoredOuterBlockConfig{}
		if err := clnt.ReadAndDecodeConfig(ctx, c, ConfigTypeFull); err != nil {
			t.Fatalf("expected bad config block to be ignored: %v", err)
		}
	})
}

func Test_EncodeDecode(t *testing.T) {
	t.Parallel()
	SkipOrFailTestIfNotRoot(t)
	slog.SetDefault(slog.New(NewContextPropagatingSlogHandler(NewTestingHandler(t))))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx := context.Background()
	clnt := GetTestClient(ctx)

	type structConfig struct {
		Config struct {
			ProfileDir string `lvm:"profile_dir"`
		} `lvm:"config"`
	}

	c := &structConfig{}
	if err := clnt.ReadAndDecodeConfig(ctx, c, ConfigTypeFull); err != nil {
		t.Fatalf("failed to get config: %v", err)
	}
	if len(c.Config.ProfileDir) == 0 {
		t.Fatalf("profile dir is empty even though that was not expected")
	}

	profileName := "lvm2go-test-encode-decode"

	testFile, err := os.OpenFile(filepath.Join(c.Config.ProfileDir, fmt.Sprintf("%s.profile", profileName)), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer func() {
		if err := testFile.Close(); err != nil {
			t.Fatalf("failed to close test file: %v", err)
		}
		if err := os.Remove(testFile.Name()); err != nil {
			t.Fatalf("failed to remove test file: %v", err)
		}
	}()

	if err := clnt.WriteAndEncodeConfig(ctx, c, testFile); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	c = &structConfig{}
	if err := clnt.ReadAndDecodeConfig(ctx, c, ConfigTypeFull, Profile(profileName)); err == nil {
		t.Fatalf("expected error due no customizable profile")
	} else if !IsConfigurationSectionNotCustomizableByProfile(err) {
		t.Fatalf("expected error due no customizable profile, but got %v", err)
	}
}

func TestGetProfilePath(t *testing.T) {
	t.Parallel()
	SkipOrFailTestIfNotRoot(t)
	slog.SetDefault(slog.New(NewContextPropagatingSlogHandler(NewTestingHandler(t))))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx := context.Background()
	clnt := GetTestClient(ctx)

	profileDir, err := clnt.GetProfileDirectory(ctx)
	if err != nil {
		t.Fatalf("failed to get profile directory: %v", err)
	} else if len(profileDir) == 0 {
		t.Fatalf("profile dir is empty even though that was not expected")
	}

	testCases := []struct {
		name     string
		profile  Profile
		expected string
		err      error
	}{
		{
			name: "empty",
			err:  ErrProfileNameEmpty,
		},
		{
			name:     "test",
			profile:  Profile("test"),
			expected: filepath.Join(profileDir, fmt.Sprintf("test%s", LVMProfileExtension)),
			err:      nil,
		},
		{
			name:     "test.profile",
			profile:  Profile("test.profile"),
			expected: filepath.Join(profileDir, fmt.Sprintf("test%s", LVMProfileExtension)),
			err:      nil,
		},
		{
			name:     "test (with valid directory)",
			profile:  Profile(filepath.Join(profileDir, "test")),
			expected: filepath.Join(profileDir, fmt.Sprintf("test%s", LVMProfileExtension)),
			err:      nil,
		},
		{
			name:     "test.profile (with valid directory)",
			profile:  Profile(filepath.Join(profileDir, "test.profile")),
			expected: filepath.Join(profileDir, fmt.Sprintf("test%s", LVMProfileExtension)),
			err:      nil,
		},
		{
			name:    "test (with invalid directory)",
			profile: Profile(filepath.Join("/bla", "test")),
			err:     fmt.Errorf("unexpected profile directory"),
		},
		{
			name:    "test.profile (with invalid directory)",
			profile: Profile(filepath.Join("/bla", "test.profile")),
			err:     fmt.Errorf("unexpected profile directory"),
		},
		{
			name:    "only folder",
			err:     fmt.Errorf("unexpected profile directory"),
			profile: Profile(profileDir),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path, err := clnt.GetProfilePath(ctx, tc.profile)
			if tc.err != nil && !strings.Contains(err.Error(), tc.err.Error()) {
				t.Fatalf("expected error %q, got %q", tc.err, err)
			}
			if path != tc.expected {
				t.Fatalf("expected path %s, got %s", tc.expected, path)
			}
		})
	}
}

func TestProfile(t *testing.T) {
	t.Parallel()
	SkipOrFailTestIfNotRoot(t)
	slog.SetDefault(slog.New(NewContextPropagatingSlogHandler(NewTestingHandler(t))))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx := context.Background()
	clnt := GetTestClient(ctx)

	profile := Profile("lvm2go-test-profile")

	type structConfig struct {
		Config struct {
			ProfileDir string `lvm:"profile_dir"`
		} `lvm:"config"`
	}

	c := &structConfig{}

	testFile, err := clnt.CreateProfile(ctx, c, profile)
	if err != nil {
		t.Fatalf("failed to create profile: %v", err)
	}
	defer func() {
		if err := clnt.RemoveProfile(ctx, profile); err != nil {
			t.Fatalf("failed to remove profile: %v", err)
		}
	}()

	if len(testFile) == 0 {
		t.Fatalf("profile dir is empty even though that was not expected")
	}

	err = clnt.ReadAndDecodeConfig(ctx, c, ConfigTypeFull, profile)
	if !IsConfigurationSectionNotCustomizableByProfile(err) {
		t.Fatalf("expected error due no customizable profile, but got %v", err)
	}
}

//go:embed testdata/lvm.conf
var testFile []byte

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkDecode(b)
	}
}

func benchmarkDecode(b *testing.B) {
	for range b.N {
		b.StopTimer()
		decoder := NewLexingConfigDecoder(bytes.NewReader(testFile))
		b.StartTimer()
		cfg := struct {
			Config struct {
				ProfileDir string `lvm:"profile_dir"`
			} `lvm:"config"`
		}{}
		if err := decoder.Decode(&cfg); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		if cfg.Config.ProfileDir != "/my/custom/profile_dir" {
			b.Fatalf("unexpected value: %s", cfg.Config.ProfileDir)
		}
	}
}

func TestUpdateGlobalConfig(t *testing.T) {
	LVMGlobalConfiguration = filepath.Join(t.TempDir(), "lvm.conf")
	if err := os.WriteFile(LVMGlobalConfiguration, testFile, 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	control := func() *bytes.Buffer {
		data, err := os.ReadFile(LVMGlobalConfiguration)
		if err != nil {
			t.Fatalf("failed to read test file: %v", err)
		}
		return bytes.NewBuffer(data)
	}

	clnt := GetTestClient(context.Background())

	cfg := struct {
		Config struct {
			AbortOnErrors int64  `lvm:"abort_on_errors"`
			ProfileDir    string `lvm:"profile_dir"`
		} `lvm:"config"`
	}{}

	cfg.Config.ProfileDir = "mynewprofiledir"

	if err := clnt.UpdateGlobalConfig(context.Background(), &cfg); err != nil {
		t.Fatalf("failed to update global config: %v", err)
	}

	containsFieldNewlySet := bytes.Contains(control().Bytes(), []byte(fmt.Sprintf(
		"abort_on_errors = %d\n",
		cfg.Config.AbortOnErrors,
	)))
	if !containsFieldNewlySet {
		println(control().String())
		t.Fatalf("expected field to be set, but it was not")
	}

	containsModifiedField := bytes.Contains(control().Bytes(), []byte(fmt.Sprintf(
		"profile_dir = %q\n",
		cfg.Config.ProfileDir,
	)))

	if !containsModifiedField {
		t.Fatalf("expected field to be modified, but it was not")
	}

	cfg.Config.ProfileDir = "mynewprofiledir2"

	if err := clnt.UpdateGlobalConfig(context.Background(), &cfg); err != nil {
		t.Fatalf("failed to update global config: %v", err)
	}

	if containsModifiedField = bytes.Contains(control().Bytes(), []byte(fmt.Sprintf(
		"profile_dir = %q\n",
		"mynewprofiledir2",
	))); !containsModifiedField {
		t.Fatalf("expected field to be modified, but it was not")
	}

	cfg.Config.ProfileDir = "mynewprofiledir3"

	if err := clnt.UpdateGlobalConfig(context.Background(), &cfg); err != nil {
		t.Fatalf("failed to update global config: %v", err)
	}

	if containsModifiedField = bytes.Contains(control().Bytes(), []byte(fmt.Sprintf(
		"profile_dir = %q\n",
		"mynewprofiledir3",
	))); !containsModifiedField {
		t.Fatalf("expected field to be modified, but it was not")
	}
}
