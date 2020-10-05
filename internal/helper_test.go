package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_firstToLower(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "consecutive caps",
			in:   "ID",
			want: "id",
		},
		{
			name: "single char",
			in:   "V",
			want: "v",
		},
		{
			name: "Upper camel",
			in:   "FirstName",
			want: "firstName",
		},
		{
			name: "Single lower prefix",
			in:   "vBar",
			want: "vBar",
		},
		{
			name: "Upper prefix",
			in:   "KTex",
			want: "kTex",
		},
		{
			name: "Unicode prefix",
			in:   "ÜberVar",
			want: "überVar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := firstToLower(tt.in)
			assert.Equal(t, tt.want, out)
		})
	}
}
