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
		{"* * * * *", true, nil},
		{"* * * aug *", true, nil},
		{"* * * * mon", true, nil},
		{"* * * * 0", true, nil},
		{"0 * * * *", true, nil},
		{"* 0 * * *", true, nil},
		{"* 3-23/100 * * *", true, nil},
		{"* 3-23/3 * * *", true, nil},
		{"1,2,3,4,5 * * * *", true, nil},
		{"1-59 * * * *", true, nil},
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

func TestDoFromMachineTimeComponent(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err error
	}{
		// It should handle ranges
		{"1-4 * * * *", "on minutes 1 through 4", nil},
		{"* 1-4 * * *", "every minute past the hours of 1 through 4", nil},
		// It should handle lists
		{"1,2,5 * * * *", "on minutes 1,2 and 5", nil},
		{"* 1,2,5 * * *", "every minute past the hours of 1,2 and 5", nil},
		// It should handle singular values
		{"1 * * * *", "at minute 1", nil},
		{"* 2 * * *", "every minute past 2", nil},
		// It should handle step values
		{"*/2 * * * *", "every 2 minutes", nil},
		{"* */6 * * *", "every minute past every 6 hours", nil},
	}
	for i, tt := range tests {
		cronP := NewCron()
		t.Run(fmt.Sprintf("Case %d: %v", i, tt.in), func(t *testing.T) {
			got, err := cronP.DoFromMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` \n; want `%s` \n; got  `%s`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` \n; want `%t` \n; got  `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestDoFromMachineDayComponent(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err error
	}{
		// It should handle ranges
		{"* * 1-25 * *", "every minute on the 1st through the 25th", nil},
		{"* * 1-25 * 1-3", "every minute on the 1st through the 25th and on Monday through Wednesday", nil},
		// It should handle list
		{"* * 1,2,3,25 * *", "every minute on the 1st, 2nd, 3rd and the 25th", nil},
		{"* * 1,2,3,25 * 1,5,7", "every minute on the 1st, 2nd, 3rd and the 25th and on Mondays, Fridays, and Sundays", nil},
		// It should handle singular values
		{"* * 31 * *", "every minute on the 31st", nil},
		{"* * 7 * 0", "every minute on the 7th and on Sundays", nil},
		// It should handle step values
		{"* * */2 * *", "every minute on every 2nd day of the month", nil},
	}
	for i, tt := range tests {
		cronP := NewCron()
		t.Run(fmt.Sprintf("Case %d: %v", i, tt.in), func(t *testing.T) {
			got, err := cronP.DoFromMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` \n; want `%s` \n; got  `%s`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` \n; want `%t` \n; got  `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestDoFromMachine(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err error
	}{
		// It should bubble up parsing errors
		{"1 2 3", "", ErrUnparsable},
    // It should handle all stars (early exit)
		{"* * * * *", "every minute", nil},
		//
		// {"1-4 3-4 * * *", "On minutes 1 through 4 past the hours of 3 through 4", nil},
		// {"1-4 3-4 5-21 * *", "On minutes 1 through 4 past the hours of 3 through 4 on the 5th through the 21st", nil},
		// {"4-45 3-4 5-21 6-10 *", "On minutes 4 through 45 past the hours of 3 through 4 on the 5th through the 21st of Jun through Oct", nil},
		// {"4-45 3-4 5-21 6-10 4-7", "On minutes 4 through 45 past the hours of 3 through 4 on the 5th through the 21st and on Thu through Sun of Jun through Oct", nil},
		// {"4-45 3-4 * 6-10 4-7", "On minutes 4 through 45 past the hours of 3 through 4 on Thu through Sun of Jun through Oct", nil},
		// // It should handle steps
		// {"*/18 * * * *", "Every 18th minute", nil},
		// {"*/18 */3 * * *", "Every 18th minute of every 3rd hour", nil},
		// {"* 1-4 * * *", "every minute of every hour from 1 through 4", nil},
		// // {"* * */3 * *", "every minute on the 3rd day of every month", nil},
		// // It should handle steps
		// {"*/4 * * * *", "every 4th minute", nil},
		// // It should lists
		// // @TODO                              this comma here could be removed
		// {"1,3,7 * * * *", "At minute 1, 3, and, 7", nil},
	}

	for i, tt := range tests {
		cronP := NewCron()
		t.Run(fmt.Sprintf("Case %d: %v", i, tt.in), func(t *testing.T) {
			got, err := cronP.DoFromMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` \n; want `%s` \n; got  `%s`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` \n; want `%t` \n; got  `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}
