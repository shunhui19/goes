package lib

import "testing"

func TestUcFirst(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "lower-http", args: args{str: "http"}, want: "Http"},
		{name: "Upper-Http", args: args{str: "Http"}, want: "Http"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UcFirst(tt.args.str)
			if got != tt.want {
				t.Errorf("UcFirst() = %v, want %v", got, tt.want)
			}
			t.Logf("UcFirst() = %v, want %v", got, tt.want)
		})
	}
}
