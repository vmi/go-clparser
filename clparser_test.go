// clparser is a simple command-line parser.

package clparser

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		commandLine            string
		enableBackslashEscapes bool
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"empty command line", args{"", false}, []string{}, false},
		{"1 arg", args{"test", false}, []string{"test"}, false},
		{"3 args", args{"test1 test2 test3", false}, []string{"test1", "test2", "test3"}, false},
		{"3 args with spaces", args{" test1  test2   test3    ", false}, []string{"test1", "test2", "test3"}, false},
		{"quoted strings 1", args{"t\\e\\s\\t\\1\\''test2 '\"test3\"\"\\\"\\\\test4\"", false}, []string{"test1'test2 test3\"\\test4"}, false},
		{"quoted strings 2", args{"\"\\a\\b\\e\\E\\f\\n\\r\\t\\v\"", true}, []string{"\a\b\u001b\u001b\f\n\r\t\v"}, false},
		{"error string 1", args{"\\", false}, nil, true},
		{"error string 2", args{"'", false}, nil, true},
		{"error string 3", args{"\"", false}, nil, true},
		{"error string 4", args{"\"\\", false}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCLParser().BackslashEscapes(tt.args.enableBackslashEscapes).Parse(tt.args.commandLine)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
