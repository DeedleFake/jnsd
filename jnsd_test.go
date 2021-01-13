package jnsd_test

import (
	"testing"

	"github.com/DeedleFake/jnsd"
)

func TestIsNameValid(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  bool
	}{
		{
			name: "Valid",
			in:   "asimplename",
			out:  true,
		},
		{
			name: "TooShort",
			in:   "az",
			out:  false,
		},
		{
			name: "TooLong",
			in:   "averyveryveryveryveryverylongname",
			out:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := jnsd.IsNameValid(test.in)
			if got != test.out {
				t.Errorf("got: %v", got)
				t.Errorf("expected: %v", test.out)
			}
		})
	}
}
