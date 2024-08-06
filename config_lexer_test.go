package lvm2go_test

import (
	"bytes"
	_ "embed"
	"strings"
	"testing"

	"github.com/jakobmoellerdev/lvm2go"
)

//go:embed testdata/lextest.conf
var lexerTest []byte

func TestConfigLexer(t *testing.T) {
	lexer := lvm2go.NewConfigLexer(bytes.NewReader(lexerTest))

	tokens, err := lexer.Lex()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stringRepresentation := strings.Contains(tokens.String(), `1:1	Comment #
1:3	CommentValue Configuration section config.
1:32	EndOfStatement

2:33	Comment #
2:35	CommentValue How LVM configuration settings are handled.
2:78	EndOfStatement

3:79	Section config
3:86	SectionStart {
4:87	EndOfStatement

4:89	Comment #
4:91	CommentValue This configuration option has an automatic default value.
4:148	EndOfStatement

5:150	Comment #
5:152	CommentValue checks = 1
5:162	EndOfStatement

7:163	EndOfStatement

7:165	Comment #
7:167	CommentValue Configuration option config/abort_on_errors.
7:211	EndOfStatement

8:213	Comment #
8:215	CommentValue Abort the LVM process if a configuration mismatch is found.
8:274	EndOfStatement

9:276	Comment #
9:278	CommentValue This configuration option has an automatic default value.
9:335	EndOfStatement

10:337	Comment #
10:339	CommentValue abort_on_errors = 0
10:358	EndOfStatement

12:359	EndOfStatement

12:361	Identifier some_field
12:372	Assignment =
12:374	Int64 1
12:376	Comment #
12:378	CommentValue This is a comment
12:395	EndOfStatement

14:396	EndOfStatement

14:398	Comment #
14:400	CommentValue Configuration option config/profile_dir.
14:440	EndOfStatement

15:442	Comment #
15:444	CommentValue Directory where LVM looks for configuration profiles.
15:497	EndOfStatement

16:499	Comment #
16:501	CommentValue This configuration option has an automatic default value.
16:558	EndOfStatement

17:560	Comment #
17:562	CommentValue profile_dir = "/etc/lvm/profile"
17:594	EndOfStatement

19:595	EndOfStatement

19:597	Identifier profile_dir
19:609	Assignment =
19:611	String /my/custom/profile_dir
19:635	EndOfStatement

20:636	SectionEnd }
-1:-1	EOF
`)
	if !stringRepresentation {
		t.Fatalf("unexpected output:\n%s", tokens.String())
	}
}
