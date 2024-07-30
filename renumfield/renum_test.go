package renumfield

import (
	"testing"
)

func Test_renum(t *testing.T) {
	type args struct {
		origText string
		fieldId  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"//", args{"", 999}, ""},
		{"//", args{"//", 999}, "//"},
		{"//", args{"// ttt", 999}, "// ttt"},
		{"//", args{" //", 999}, " //"},
		{"//", args{" // ttt", 999}, " // ttt"},
		{"//", args{"xx // ttt", 999}, "xx // ttt"},
		{"//", args{`xx="nn//bb;" // a comment`, 999}, `xx:999="nn//bb;" // a comment`},

		{";", args{";", 999}, ";"},
		{";", args{" ;", 999}, " ;"},
		{";", args{" ; ", 999}, " ; "},
		{";", args{"tt;", 999}, "tt:999;"},
		{";", args{"tt;// xx", 999}, "tt:999;// xx"},
		{";", args{"tt; // xx", 999}, "tt:999; // xx"},
		{";", args{"string tt;", 999}, "string tt:999;"},
		{";", args{"string tt; // xx", 999}, "string tt:999; // xx"},
		{";", args{`string tt = "abc;def";`, 999}, `string tt:999 = "abc;def";`},

		{"=", args{"tt=111;", 999}, "tt:999=111;"},
		{"=", args{"tt =111;", 999}, "tt:999 =111;"},
		{"=", args{"tt = 111;", 999}, "tt:999 = 111;"},
		{"=", args{"tt = 111; // xxx", 999}, "tt:999 = 111; // xxx"},

		{":", args{"tt:222;", 999}, "tt:999;"},
		{":", args{"tt :222;", 999}, "tt:999;"},
		{":", args{"tt: 222;", 999}, "tt:999;"},
		{":", args{"tt : 222;", 999}, "tt:999;"},

		{":=", args{"tt:222=xx;", 999}, "tt:999=xx;"},
		{":=", args{"tt:222 =xx;", 999}, "tt:999 =xx;"},
		{":=", args{"tt:222= xx;", 999}, "tt:999= xx;"},
		{":=", args{"tt:222 = xx;", 999}, "tt:999 = xx;"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := RenumLine([]byte(tt.args.origText), tt.args.fieldId); string(got) != tt.want {
				t.Errorf("renum([%s]) =\n [%s], want\n [%s]", tt.args.origText, string(got), tt.want)
			}
		})
	}
}
