package lib

import (
	"fmt"
	"testing"
)

func Test_textColor(t *testing.T) {
	type args struct {
		color int
		str   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"black", args{30, "black"}, "black"},
		{"red", args{31, "red"}, "red"},
		{"green", args{32, "green"}, "green"},
		{"yellow", args{33, "yellow"}, "yellow"},
		{"blue", args{34, "blue"}, "blue"},
		{"magenta", args{35, "magenta"}, "magenta"},
		{"cyan", args{36, "cyan"}, "cyan"},
		{"white", args{37, "white"}, "white"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := textColor(tt.args.color, tt.args.str)
			if got != fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", tt.args.color, tt.want) {
				t.Errorf("textColor() = %v, want %v", got, tt.want)
			}
			t.Logf("textColor() = %v, want %v\n", got, tt.want)
		})
	}
}
