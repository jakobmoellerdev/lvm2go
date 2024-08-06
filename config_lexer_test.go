package lvm2go_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/jakobmoellerdev/lvm2go"
)

//go:embed testdata/lextest.conf
var lexerTest []byte

//go:embed testdata/lextest.output
var lexTestOutput string

func TestConfigLexer(t *testing.T) {
	lexer := lvm2go.NewBufferedConfigLexer(bytes.NewReader(lexerTest))

	tokens, err := lexer.Lex()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !(tokens.String() == lexTestOutput) {
		t.Fatalf("unexpected output:\n%s", tokens.String())
	}

	data, err := lvm2go.ConfigTokensToBytes(tokens)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != string(lexerTest) {
		t.Fatalf("unexpected output:\n%s", data)
	}
}

func TestNewLexingConfigDecoder(t *testing.T) {

	t.Run("structured", func(t *testing.T) {
		decoder := lvm2go.NewLexingConfigDecoder(bytes.NewReader(lexerTest))
		cfg := struct {
			Config struct {
				SomeField  int64  `lvm:"some_field"`
				ProfileDir string `lvm:"profile_dir"`
			} `lvm:"config"`
		}{}

		if err := decoder.Decode(&cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Config.SomeField != 1 {
			t.Fatalf("unexpected value: %d", cfg.Config.SomeField)
		}
		if cfg.Config.ProfileDir != "/my/custom/profile_dir" {
			t.Fatalf("unexpected value: %s", cfg.Config.ProfileDir)
		}
	})

	t.Run("unstructured", func(t *testing.T) {
		decoder := lvm2go.NewLexingConfigDecoder(bytes.NewReader(lexerTest))
		cfg := map[string]any{}

		if err := decoder.Decode(&cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg["config/some_field"].(int64) != 1 {
			t.Fatalf("unexpected value: %d", cfg["config/some_field"])
		}
		if cfg["config/profile_dir"].(string) != "/my/custom/profile_dir" {
			t.Fatalf("unexpected value: %s", cfg["config/profile_dir"])
		}
	})
}

func TestNewLexingConfigEncoder(t *testing.T) {
	t.Run("structured", func(t *testing.T) {
		cfg := struct {
			Config struct {
				SomeField  int64  `lvm:"some_field"`
				ProfileDir string `lvm:"profile_dir"`
			} `lvm:"config"`
		}{}

		cfg.Config.SomeField = 1
		cfg.Config.ProfileDir = "/my/custom/profile_dir"

		testBuffer := &bytes.Buffer{}
		encoder := lvm2go.NewLexingConfigEncoder(testBuffer)

		if err := encoder.Encode(&cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if testBuffer.String() != `config {
	some_field = 1
	profile_dir = "/my/custom/profile_dir"
}
` {
			t.Fatalf("unexpected output:\n%s", testBuffer.String())
		}
	})

	t.Run("unstructured", func(t *testing.T) {
		cfg := map[string]any{}

		cfg["config/some_field"] = int64(1)
		cfg["config/profile_dir"] = "/my/custom/profile_dir"

		testBuffer := &bytes.Buffer{}
		encoder := lvm2go.NewLexingConfigEncoder(testBuffer)

		if err := encoder.Encode(&cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if testBuffer.String() != `config {
	profile_dir = "/my/custom/profile_dir"
	some_field = 1
}
` {
			t.Fatalf("unexpected output:\n%s", testBuffer.String())
		}
	})
}
