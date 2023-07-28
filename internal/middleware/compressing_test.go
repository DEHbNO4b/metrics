package middleware

import "testing"

func Test_isNeedToCompress(t *testing.T) {

	tests := []struct {
		name string
		s    []string
		want bool
	}{{
		name: "positive test#1",
		s:    []string{"text/html"},
		want: true,
	},
		{
			name: "positive test#2",
			s:    []string{"application/json"},
			want: true,
		},
		{
			name: "empty string",
			s:    []string{},
			want: false,
		},
		{
			name: "number",
			s:    []string{"343"},
			want: false,
		},
		{
			name: "negative test#2",
			s:    []string{"sometext"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNeedToCompress(tt.s); got != tt.want {
				t.Errorf("isNeedToCompress() = %v, want %v", got, tt.want)
			}
		})
	}
}
