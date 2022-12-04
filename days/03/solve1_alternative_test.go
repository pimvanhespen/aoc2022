package main

import (
	"bytes"
	"fmt"
	"testing"
)

var longstring = `LdnrrLnhRdLLmLDRPvmdQnJDJWNqcCqZJZqfFqfcfzcq
vPTbfWggzvGVqjsVqV
dDcJHZcZHmMFQQMshsjcRqVChjNtqh
`

func Test_solve1_stdlib(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Simple string",
			args: args{
				b: []byte("aaaabbba\n"),
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Double string",
			args: args{
				b: []byte("aaaabbba\naaaabbba\n"),
			},
			want:    2,
			wantErr: false,
		},

		{
			name: "long string",
			args: args{
				b: []byte(longstring),
			},
			want:    55,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := solve1_noAllocs(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("solve1_noAllocs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("solve1_noAllocs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSth(t *testing.T) {
	s := []byte(`vPTbfWggzvGVqjsVqV`)

	for _, r := range s {
		fmt.Printf("%3d ", r)
	}
	fmt.Println()
	fmt.Println(s[:9])
	fmt.Println(s[9:])

	for _, r := range s[:9] {
		search := bytes.IndexByte(s[9:], r)
		fmt.Printf("%c %3d --> %d\n", r, r, search)
	}

}
