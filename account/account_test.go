package account

import (
	"testing"
)

func TestPassword_Compare(t *testing.T) {
	tests := []struct {
		name string
		p    Password
		args Password
		want bool
	}{
		{
			p:    "abcdefgh",
			args: Password("abcdefgh"),
			want: true,
		}, {
			p:    "abcdefgh",
			args: Password("bcdefgh"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Compare(tt.args); got != tt.want {
				t.Errorf("Password.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkPassword_Hash(b *testing.B) {
	p := Password("djfajfalj")
	for i := 0; i < b.N; i++ {
		_ = p.Hash()
	}
}
