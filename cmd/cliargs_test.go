package cmd

import (
	"reflect"
	"testing"
)

func TestIsFlag(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"foo", false},
		{"-foo", true},
		{"--foo", false},
	}
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := IsFlag(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
		})
	}
}

func TestIsOption(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"foo", false},
		{"-foo", false},
		{"--foo", true},
	}
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := IsOption(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
		})
	}
}

func TestIsPositional(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"foo", true},
		{"-foo", false},
		{"--foo", false},
	}
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := IsPositional(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
		})
	}
}

func TestParseCliArgs(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		out  CliArgs
	}{
		// human foo baz
		{
			"Positionals Happy Path",
			[]string{"foo", "baz"},
			CliArgs{
				Flags:       map[string]bool{},
				Options:     map[string]string{},
				Positionals: []string{"foo", "baz"},
			},
		},
		// human --foo=bar
		{
			"Options Happy Path",
			[]string{"--foo=bar"},
			CliArgs{
				Flags:       map[string]bool{},
				Options:     map[string]string{"foo": "bar"},
				Positionals: []string{},
			},
		},
		// human --foo
		{
			"Options without values",
			[]string{"--foo"},
			CliArgs{
				Flags:       map[string]bool{},
				Options:     map[string]string{"foo": ""},
				Positionals: []string{},
			},
		},
		// human --foo bar --baz
		{
			"Options swallow next positional",
			[]string{"--foo", "bar", "--baz"},
			CliArgs{
				Flags:       map[string]bool{},
				Options:     map[string]string{"foo": "bar", "baz": ""},
				Positionals: []string{},
			},
		},
		// human --foo bar --foo baz
		{
			"Options last option wins",
			[]string{"--foo", "bar", "--foo", "baz"},
			CliArgs{
				Flags:       map[string]bool{},
				Options:     map[string]string{"foo": "baz"},
				Positionals: []string{},
			},
		},
		// human -f
		{
			"Flags Happy Path",
			[]string{"-f"},
			CliArgs{
				Flags:       map[string]bool{"f": true},
				Options:     map[string]string{},
				Positionals: []string{},
			},
		},
		// human -b -a -r
		{
			"Flags Happy Path",
			[]string{"-b", "-a", "-r"},
			CliArgs{
				Flags:       map[string]bool{"b": true, "a": true, "r": true},
				Options:     map[string]string{},
				Positionals: []string{},
			},
		},
		// human -bar
		{
			"Flags handle shorthand",
			[]string{"-bar"},
			CliArgs{
				Flags:       map[string]bool{"b": true, "a": true, "r": true},
				Options:     map[string]string{},
				Positionals: []string{},
			},
		},
		// human -foo
		// adding this for usage clarity, the fact that the implementation uses a
		// map takes care of this however wanted to bubble up the behavior
		{
			"Flags handle repetitive flags",
			[]string{"-foo"},
			CliArgs{
				Flags:       map[string]bool{"f": true, "o": true},
				Options:     map[string]string{},
				Positionals: []string{},
			},
		},
		{
			"Kitchen Sink",
			[]string{"first", "-foo", "-b", "-a", "-r", "--opt", "--foo=bar", "baz", "--long", "wat", "last"},
			CliArgs{
				Flags:       map[string]bool{"f": true, "o": true, "b": true, "a": true, "r": true},
				Options:     map[string]string{"foo": "bar", "opt": "", "long": "wat"},
				Positionals: []string{"first", "baz", "last"},
			},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCliArgs(tt.in)
			conds := []bool{
				reflect.DeepEqual(got.Flags, tt.out.Flags),
				reflect.DeepEqual(got.Options, tt.out.Options),
				reflect.DeepEqual(got.Positionals, tt.out.Positionals),
			}
			for _, c := range conds {
				if c == false {
					t.Errorf("\nCase %d: \nGiven = `%v` ; \nwant `%v` ; \ngot `%v`", i, tt.in, tt.out, got)
				}
			}
		})
	}
}
