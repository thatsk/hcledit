package editor

import (
	"bytes"
	"testing"
)

func TestAttributeAppend(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		address string
		value   string
		newline bool
		ok      bool
		want    string
	}{
		{
			name: "simple top level attribute",
			src: `
a0 = v0
`,
			address: "a1",
			value:   "v1",
			newline: false,
			ok:      true,
			want: `
a0 = v0
a1 = v1
`,
		},
		{
			name: "attribute in block",
			src: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
			address: "b1.l1.a2",
			value:   "v2",
			newline: false,
			ok:      true,
			want: `
a0 = v0
b1 "l1" {
  a1 = v1
  a2 = v2
}
`,
		},
		{
			name: "attribute in block (with newline)",
			src: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
			address: "b1.l1.a2",
			value:   "v2",
			newline: true,
			ok:      true,
			want: `
a0 = v0
b1 "l1" {
  a1 = v1

  a2 = v2
}
`,
		},
		{
			name: "block not found (noop)",
			src: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
			address: "b2.l1.a1",
			value:   "v2",
			newline: false,
			ok:      true,
			want: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
		},
		{
			name: "attribute already exists (error)",
			src: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
			address: "b1.l1.a1",
			value:   "v2",
			newline: false,
			ok:      false,
			want:    ``,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inStream := bytes.NewBufferString(tc.src)
			outStream := new(bytes.Buffer)
			err := AppendAttribute(inStream, outStream, "test", tc.address, tc.value, tc.newline)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			got := outStream.String()
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, outStream: \n%s", got)
			}

			if got != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}
