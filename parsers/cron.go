package parsers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

var ErrBadRange error = errors.New("Ranges can only contain 2 numbers a start and and end. One of the ranges received was incorrectly formated")
var ErrHangingRangeList error = errors.New("Ranges and/or Lists is hanging meaning that either a comma or a hyphen isn't followed by a number: ie: N- or N,")
var ErrBadMinuteField error = errors.New("Bad minute field. The value received is not a number between 1-59")
var ErrBadHourField error = errors.New("Bad hour field. The value received is not a number between 1-23")
var ErrBadDomField error = errors.New("Bad day of month field. The value received is not a number between 1-31")
var ErrBadMonthField error = errors.New("Bad month field. The value received is not a number between 1-12")
var ErrBadDowField error = errors.New("Bad day of week field. The value received is not a number between 0-7")
var ErrBadRangeStep error = errors.New("Bad step provided for range syntax. This means that the step is either smaller or larger than the lower and upper bounds of the time part where it was found or the step number could not be parsed or the syntax was incorrect and more than one '/' was found")

// getSliceOfNumbers parses the range and list syntaxes into slices of numbers.
// `startStop` is used when the input is an asterisk
func getSliceOfNumbers(s string, startStop string) ([]int64, error) {
	var rtn, emptyRtn []int64

	// Ensure we don't have a hanging range or step
	if s[len(s)-1] == ',' || s[len(s)-1] == '-' {
		return emptyRtn, ErrHangingRangeList
	}

	// Since both lists of numbers and ranges can be passed we'll create a
	// list of items and process only 1 item at a time; being a number or a range
	if strings.Contains(s, ",") {
		items := strings.Split(s, ",")
		for _, v := range items {
			r, err := getSliceOfNumbers(v, startStop)
			if err != nil {
				return emptyRtn, err
			}
			rtn = append(rtn, r...)
		}
		return rtn, nil
	}

	// Handle Ranges Steps (*/[step] or n-n/[step])
	//
	// A range can include a "step increment" by using the forward slash and can
	// be used after asterisks or a hyphenated input.
	//
	// In order for a step to be valid it should be reater than 0
	step := int64(1)
	if strings.Contains(s, "/") {
		chunks := strings.Split(s, "/")
		rawStep := chunks[len(chunks)-1]

		// Ensure step is valid
		gotstep, err := strconv.ParseInt(rawStep, 10, 64)
		if len(chunks) > 2 || gotstep == 0 || err != nil {
			return emptyRtn, ErrBadRangeStep
		}

		step = gotstep
		s = chunks[0]

		// When an asterisk is used we can short-circuit since we know the start
		// and stop of the range given which "time part" the asterisk appears in
		if s == "*" {
			return getSliceOfNumbers(startStop+"/"+rawStep, startStop)
		}

	}

	// Handle a range (hyphenated input)
	if strings.Contains(s, "-") {
		nums := []int64{}

		for _, v := range strings.SplitN(s, "-", 2) {
			// Only a range between 2 numbers is allowed
			if strings.Contains(v, "-") {
				return emptyRtn, ErrBadRange
			}

			n, err := getSliceOfNumbers(v, startStop)
			if err != nil {
				return emptyRtn, err
			}
			nums = append(nums, n[0])
		}

		// Given that we're dealing with a range in this branch we can make the
		// slice with what we know now
		cur := nums[0]
		for cur <= nums[1] {
			rtn = append(rtn, cur)
			cur = cur + step
		}

		return rtn, nil
	}

	// At this point `s` should be a number string so we should be able to parse it
	{
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return emptyRtn, err
		}

		rtn = append(rtn, n)
	}

	return rtn, nil

}

func invalidMinute(a int64) bool {
	return a < 0 || a > 59
}

func invalidHour(a int64) bool {
	return a < 0 || a > 23
}

func invalidDom(a int64) bool {
	return a < 1 || a > 31
}

func invalidMonth(a int64) bool {
	return a < 1 || a > 12
}

func invalidDow(a int64) bool {
	return a < 0 || a > 7
}

type parsedOutput struct {
	minutes []int64
	hours   []int64
	dom     []int64
	month   []int64
	dow     []int64
}

// outputShape stores the components (english words) that are used to build the
// human readable output
type outputShape struct {
	minutes string
	hours   string
	dom     string
	month   string
	dow     string
}

type Cron struct {
	monthNames [12]string
	dowNames   [7]string
}

func NewCron() *Cron {
	return &Cron{
		[12]string{
			"jan",
			"feb",
			"mar",
			"apr",
			"may",
			"jun",
			"jul",
			"aug",
			"sep",
			"oct",
			"nov",
			"dec",
		},
		[7]string{
			"mon",
			"tue",
			"wed",
			"thu",
			"fri",
			"sat",
			"sun",
		},
	}
}

// parseInputOrError will return
func (c *Cron) parseInputOrError(input string) (parsedOutput, error) {
	var rtn parsedOutput

	// all five fields must be present.
	rawParts := strings.Split(input, " ")
	if len(rawParts) != 5 {
		return rtn, ErrUnparsable
	}

	// All fields allow for numbers, letters, the hyphen, asterisk,
	// forward-slashes and the comma so we'll go broad with the expression just
	// to knock out obvious unknown stuff
	//
	// In particular the key takeaway is that the first character isn't optional
	// ie: the "match 1 token" ( `{1}` )
	//
	// We'll check each field individually later on to keep the regex
	// quickly understandable
	r := regexp.MustCompile(`(?i)^[0-9a-z\*/]{1}[0-9,\-*a-z/]*$`)
	for _, v := range rawParts {
		if !r.MatchString(v) {
			return rtn, ErrUnparsable
		}
	}

	// minute, hour, dom, month, dow = rawParts
	rawMinute := rawParts[0]
	rawHour := rawParts[1]
	rawDom := rawParts[2]
	rawMonth := rawParts[3]
	rawDow := rawParts[4]

	// Check fields that don't allow letters
	for _, f := range []string{rawMinute, rawHour, rawDom} {
		if notok, _ := regexp.MatchString(`(?i)[a-z]`, f); notok {
			return rtn, ErrUnparsable
		}
	}

	// Okay now we can check the individual fields
	//
	// "*" means `first-last` so we'll convert those now
	if rawMinute == "*" {
		rawMinute = "0-59"
	}

	if rawHour == "*" {
		rawHour = "0-23"
	}

	if rawDom == "*" {
		rawDom = "1-31"
	}

	if rawMonth == "*" {
		rawMonth = "1-12"
	}

	if rawDow == "*" {
		rawDow = "0-7"
	}

	// Deal with Lists and Ranges
	//
	// The idea here is to take the cron syntax that represents a range and/or a
	// list and parse it into a slice so we can work with numbers instead
	minutes, err := getSliceOfNumbers(rawMinute, "0-59")
	if err != nil {
		return rtn, err
	}

	hours, err := getSliceOfNumbers(rawHour, "0-23")
	if err != nil {
		return rtn, err
	}

	dom, err := getSliceOfNumbers(rawDom, "1-31")
	if err != nil {
		return rtn, err
	}

	// Since these fields allow abbreviations we need to also check for that
	month, err := func() ([]int64, error) {
		var m []int64
		m, e := getSliceOfNumbers(rawMonth, "1-12")
		_, intParseError := e.(*strconv.NumError)
		if e != nil && intParseError {
			for k, v := range c.monthNames {
				if rawMonth == v {
					m = append(m, int64(k)+1)
					return m, nil
				}
			}
			return m, ErrUnparsable
		}
		return m, e
	}()
	if err != nil {
		return rtn, err
	}

	dow, err := func() ([]int64, error) {
		var d []int64
		d, e := getSliceOfNumbers(rawDow, "0-7")
		_, intParseError := e.(*strconv.NumError)
		if e != nil && intParseError {
			for k, v := range c.dowNames {
				if rawDow == v {
					d = append(d, int64(k)+1)
					return d, nil
				}
			}
			return d, ErrUnparsable
		}
		return d, e
	}()
	if err != nil {
		return rtn, err
	}

	// Okay so now we can validate the actual values we parsed
	for _, m := range minutes {
		if invalidMinute(m) {
			return rtn, ErrBadMinuteField
		}
	}

	for _, h := range hours {
		if invalidHour(h) {
			return rtn, ErrBadHourField
		}
	}

	for _, d := range dom {
		if invalidDom(d) {
			return rtn, ErrBadDomField
		}
	}

	for _, m := range month {
		if invalidMonth(m) {
			return rtn, ErrBadMonthField
		}
	}

	for _, w := range dow {
		if invalidDow(w) {
			return rtn, ErrBadDowField
		}
	}

	rtn.minutes = minutes
	rtn.hours = hours
	rtn.dom = dom
	rtn.month = month
	rtn.dow = dow

	return rtn, nil
}

// CanParseFromMachine will determine if input will work for us.
//
// Need to be:
// - 5 fields
// - fields are separated by a space
//
//    ```
//    field         allowed values (all accept * as well)
//    -----         --------------
//    minute        0-59
//    hour          0-23
//    day of month  1-31
//    month         1-12 (or names, see below)
//    day of week   0-7  (0 or 7 is Sun, or use names)
//    ```
func (c *Cron) CanParseFromMachine(input string) (bool, error) {
	if _, err := c.parseInputOrError(input); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Cron) DoIntoMachine(string) (string, error) {
	return ErrNotYetImplemented.Error(), nil
}

func (c *Cron) DoFromMachine(input string) (string, error) {
	parsed, err := c.parseInputOrError(input)
	if err != nil {
		return "", err
	}

	output := outputShape{}

	// Alias some facts
	everyMinute := len(parsed.minutes) == 60
	everyHour := len(parsed.hours) == 24
	everyDom := len(parsed.dom) == 31
	everyMonth := len(parsed.month) == 12
	everyDow := len(parsed.dow) == 8

	// Logic for human formatting looks like:
	//
	// Each "component" is joined with the joiner token "of" since this has the
	// most flexibility from an english standpoint
	joinerToken := "of"

	// The actual english values for each component can be mapped to the their
	// numeric counterpart with the following exceptions:
	//
	// if any value is set to *
	//   then the output value should be "every $value"

	if everyMinute {
		output.minutes = "every minute"
	}

	if everyHour {
		output.hours = "every hour"
	}

	if everyDom {
		output.dom = "every day"
	}

	if everyMonth {
		output.month = "every month"
	}

	if everyDow {
		output.dow = "every day of the week"
	}

	// If every component is set to "every", then we can reduce the noise in the
	// output by allow the reader to infer the other time components
	if func() bool {
		for _, c := range []bool{everyMinute, everyHour, everyDom, everyMonth, everyDow} {
			if c != true {
				return false
			}
		}
		return true
	}() {
		fmt.Println("OUTPUT:", "every minute")
		return "every minute", nil
	}

	// spew.Dump(parsed)
	spew.Dump(output)
	str := strings.Join([]string{
		output.minutes,
		output.hours,
		output.dom,
		output.month,
		output.dow,
	}, " "+joinerToken+" ",
	)
	fmt.Println("str", str)

	// allStars := func() bool {
	// 	for _, r := range rawParts {
	// 		if r != "*" {
	// 			return false
	// 		}
	// 	}
	// 	return true
	// }()
	// if allStars == true {
	// 	return true, nil
	// }
	// spew.Dump(allStars)

	//// if both month and dom are set to *:
	//// 		then month and dom become "daily"
	//// else:
	//// 		for month:
	//// 			{month} {dom}
	////
	//// if day of week is not present:
	//// 		[{month}] [{day}] at {hour minute}
	//// else:
	//// 		[{month}] [{day}] at {hour minute} and {dow}

	//// fmt.Println("OUTPUT")
	//// fmt.Println("for input")
	//// fmt.Println(input)
	//// fmt.Println("minutes")
	//// spew.Dump(minutes)
	//// fmt.Println("hours")
	//// spew.Dump(hours)
	//// fmt.Println("dom")
	//// spew.Dump(dom)
	//// fmt.Println("month")
	//// spew.Dump(month)
	//// fmt.Println("dow")
	//// spew.Dump(dow)
	//// fmt.Println("=================")
	return ErrNotYetImplemented.Error(), nil
}
