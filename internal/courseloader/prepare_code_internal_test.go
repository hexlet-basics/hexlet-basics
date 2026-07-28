package courseloader

import "testing"

// TestPrepareCode covers the BEGIN/END solution-stripping in isolation — the one
// piece of non-mechanical logic in the parser. It is an in-package test because
// prepareCode is unexported.
func TestPrepareCode(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no markers yields empty template",
			in:   "console.log(\"Hello, World!\");\n",
			want: "",
		},
		{
			// The begin group's trailing \s* swallows the newline + next-line
			// indent, so a blanked block leaves that indent as its own line — this
			// is exactly what the legacy Ruby regex produces, preserved for parity.
			name: "single block is blanked and relabeled",
			in:   "def solution\n  # BEGIN\n  42\n  # END\nend\n",
			want: "def solution\n  # BEGIN (write your solution here)\n  \n  # END\nend\n",
		},
		{
			name: "two blocks each collapse independently",
			in:   "# BEGIN\na\n# END\nx = 1\n# BEGIN\nb\n# END\n",
			want: "# BEGIN (write your solution here)\n\n# END\nx = 1\n# BEGIN (write your solution here)\n\n# END\n",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := prepareCode(tc.in)
			if got != tc.want {
				t.Fatalf("prepareCode(%q)\n got: %q\nwant: %q", tc.in, got, tc.want)
			}
		})
	}
}
