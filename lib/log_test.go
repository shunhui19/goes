package lib

import (
	"fmt"
	"testing"
)

func TestLogPrint(t *testing.T) {
	type args struct {
		format string
		a      []interface{}
	}
	tests := []struct {
		name string
		args args
		t    string
	}{
		{name: "Info-without-args", args: args{format: "The type of Info message"}, t: "Info"},
		{name: "Info-with-args", args: args{format: "The type of Info message with:%d, %v", a: []interface{}{15, "info"}}, t: "Info"},
		{name: "Info-with-one-args", args: args{format: "this is Info message with: %v", a: []interface{}{"hello"}}, t: "Info"},

		{name: "Debug-without-args", args: args{format: "The type of Debug message"}, t: "Debug"},
		{name: "Debug-with-args", args: args{format: "The type of Debug message with:%d, %v", a: []interface{}{15, "debug"}}, t: "Debug"},
		{name: "Debug-with-one-args", args: args{format: "this is Debug message with: %v", a: []interface{}{"hello"}}, t: "Debug"},

		{name: "Warn-without-args", args: args{format: "The type of Warn message"}, t: "Warn"},
		{name: "Warn-with-args", args: args{format: "The type of Warn message with:%d, %v", a: []interface{}{15, "Warn"}}, t: "Warn"},
		{name: "Warn-with-one-args", args: args{format: "this is Warn message with: %v", a: []interface{}{"hello"}}, t: "Warn"},

		{name: "Error-without-args", args: args{format: "The type of Error message"}, t: "Error"},
		{name: "Error-with-args", args: args{format: "The type of Error message with:%d, %v", a: []interface{}{15, "Error"}}, t: "Error"},
		{name: "Error-with-one-args", args: args{format: "this is Error message with: %v", a: []interface{}{"hello"}}, t: "Error"},

		{name: "Panic-without-args", args: args{format: "The type of Panic message"}, t: "Panic"},
		{name: "Panic-with-args", args: args{format: "The type of Panic message with:%d, %v", a: []interface{}{15, "Panic"}}, t: "Panic"},
		{name: "Panic-with-one-args", args: args{format: "this is Panic message with: %v", a: []interface{}{"hello"}}, t: "Panic"},

		{name: "Fatal-without-args", args: args{format: "The type of Fatal message"}, t: "Fatal"},
		{name: "Fatal-with-args", args: args{format: "The type of Fatal message with:%d, %v", a: []interface{}{15, "Fatal"}}, t: "Fatal"},
		{name: "Fatal-with-one-args", args: args{format: "this is Fatal message with: %v", a: []interface{}{"hello"}}, t: "Fatal"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.t {
			case "Info":
				Info(tt.args.format, tt.args.a...)
			case "Debug":
				Debug(tt.args.format, tt.args.a...)
			case "Warn":
				Warn(tt.args.format, tt.args.a...)
			case "Error":
				Error(tt.args.format, tt.args.a...)
			case "Panic":
				defer func() {
					if p := recover(); p != nil {
						fmt.Println(p)
					}
				}()
				Panic(tt.args.format, tt.args.a...)
			// note! because Fatal will call os.Exit(1), so this method will cancel
			case "Fatal":
				//Fatal(tt.args.format, tt.args.a...)
				break
			default:
				t.Logf("no match")
			}
		})
	}
}
