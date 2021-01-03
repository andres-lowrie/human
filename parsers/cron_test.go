package parsers

import (
	"fmt"
	"testing"
)

func TestCronCanParseFromMachine(t *testing.T) {
	tests := []struct {
		in  string
		out bool
		err error
	}{
		// It should fail if it doesn't have 5 fields
		{"A B C", false, ErrUnparsable},
		{"A B C D E F", false, ErrUnparsable},
		// It should fail on hanging lists and ranges
		{", B , D ,", false, ErrUnparsable},
		{"A - C - E", false, ErrUnparsable},
		// It should fail if a field has a hanging list or range
		{"0, * * * *", false, ErrHangingRangeList},
		{"59- * * * *", false, ErrHangingRangeList},
		{"* 0, * * *", false, ErrHangingRangeList},
		{"* 59- * * *", false, ErrHangingRangeList},
		{"* * 0, * *", false, ErrHangingRangeList},
		{"* * 59- * *", false, ErrHangingRangeList},
		{"* * * 0, *", false, ErrHangingRangeList},
		{"* * * 59- *", false, ErrHangingRangeList},
		{"* * * * 0,", false, ErrHangingRangeList},
		{"* * * * 59-", false, ErrHangingRangeList},
		// It should fail if an incorrect range is passed
		{"2-4-6 * * * *", false, ErrBadRange},
		{"1-2,2-4-6 * * * *", false, ErrBadRange},
		{"* 7-5-6 * * *", false, ErrBadRange},
		{"* 10-14,9-10-12 * * *", false, ErrBadRange},
		{"* * 7-5-6 * *", false, ErrBadRange},
		{"* * 7-5-6,5-9 * *", false, ErrBadRange},
		{"* * * 1-2-3 *", false, ErrBadRange},
		{"* * * 1-2-3,6-9 *", false, ErrBadRange},
		{"* * * * 1-2-3", false, ErrBadRange},
		{"* * * * 1-2-3,6-9", false, ErrBadRange},
		// Field ranges
		{"60 * * * *", false, ErrBadMinuteField},
		{"* 60 * * *", false, ErrBadHourField},
		{"* * 60 * *", false, ErrBadDomField},
		{"* * * 60 *", false, ErrBadMonthField},
		{"* * * * 60", false, ErrBadDowField},
		{"*//2 * * * *", false, ErrBadRangeStep},
		{"* *//2 * * *", false, ErrBadRangeStep},
		{"* * *//2 * *", false, ErrBadRangeStep},
		{"* * * *//2 *", false, ErrBadRangeStep},
		{"* * * * *//2", false, ErrBadRangeStep},
		{"* * */0 * *", false, ErrBadRangeStep},
		{"* * * */0 *", false, ErrBadRangeStep},
		{"59/* * * * *", false, ErrBadRangeStep},
		{"* 23/* * * *", false, ErrBadRangeStep},
		{"* * 1/* * *", false, ErrBadRangeStep},
		{"* * * 1/* *", false, ErrBadRangeStep},
		{"* * * * 7/*", false, ErrBadRangeStep},
		// Non asterisk range checking
		{"* 0-23/0 * * *", false, ErrBadRangeStep},
		// Fields with names
		// It should fail if letters are present on the non-letter fields
		{"a * * * *", false, ErrUnparsable},
		{"* a * * *", false, ErrUnparsable},
		{"* * a * *", false, ErrUnparsable},
		// It should fail if month abbrev doesn't exist
		{"* * * abc *", false, ErrUnparsable},
		{"* * * * abc", false, ErrUnparsable},
		// Positive
		{"* * * aug *", false, ErrNotYetImplemented},
		{"* * * * mon", false, ErrNotYetImplemented},
		// {"* 3-23/100 * * *", false, ErrBadRangeStep},
		// {"* 3-23/3 * * *", false, ErrBadRangeStep},
		// {"*/2 * * * *", false, errors.New("foobar")},
		// simple list {"1,2,3,4,5 * * * *", false, ErrBadRange},
		// simele range {"1-59 * * * *", false, ErrBadMinuteField},
		// These allow zero values
		// {"* * * * */0", false, ErrBadRangeStep},
		// {"*/0 * * * *", false, ErrBadRangeStep},
		// {"* */0 * * *", false, ErrBadRangeStep},
	}

	for i, tt := range tests {
		cronP := NewCron()
		t.Run(fmt.Sprintf("Case %d: %v", i, tt.in), func(t *testing.T) {
			got, err := cronP.CanParseFromMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}
