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
			name: "Empty string",
			in:   "",
			want: "",
		},
		{
			name: "Just ID",
			in:   "ID",
			want: "id",
		},
		{
			name: "ID prefix",
			in:   "IDHolder",
			want: "idHolder",
		},
		{
			name: "Single char",
			in:   "V",
			want: "v",
		},
		{
			name: "Single unicode char",
			in:   "Ä",
			want: "ä",
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
		{
			name: "ID anywhere",
			in:   "CustomerID",
			want: "customerID",
		},
		{
			name: "JSON is here",
			in:   "SomeJSON",
			want: "someJSON",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := firstToLower(tt.in)
			assert.Equal(t, tt.want, out)
		})
	}
}
