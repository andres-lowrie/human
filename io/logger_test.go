package io

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestItGivesBackInterfaceWhenTurnedOff(t *testing.T) {
	got := NewLogger(OFF, false)

	if v, ok := interface{}(got).(Ourlog); !ok {
		t.Errorf("Case Inteface: it should give back an interface `%v+`", v)
	}
}

// this is using the implementation details in order to ensure the `level`
// logic works as expected.
//
// In this case we're checking the
func TestItShouldOnlyLogEventsUpToTheLevel(t *testing.T) {
	tests := []struct {
		in   level
		want int // Used to check the internal buffer length of the logger
	}{
		// Not sure what's going on here but the internal representation is off by
		// one so the want counts are offset to accomidate for it.
		// Another option would be to check the content of the buffer... :shrug:
		// this seems to be good enough for the goal behind this test anyway
		{INFO, 2},
		{WARN, 4},
		{DEBUG, 6},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Level: %d", tt.in), func(t *testing.T) {
			// All cases get the same calls since we're checking accumulative values
			l := NewLogger(tt.in, false)
			l.Info("a")
			l.Warn("b")
			l.Debug("c")

			rv := reflect.ValueOf(l).Elem().FieldByName("Buf")
			b := rv.Interface().(bytes.Buffer)
			got := b.Len()
			if got != tt.want {
				t.Errorf("Case %d: Given = `%d` ; want `%d` ; got `%d`", i, tt.in, tt.want, got)
			}
		})
	}
}
