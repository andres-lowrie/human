package parsers

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

var ErrBadRange error = errors.New("Ranges can only contain 2 numbers a start and and end. One of the ranges received was incorrectly formated")
var ErrHangingRangeList error = errors.New("Ranges and/or Lists is hanging meaning that either a comma or a hyphen isn't followed by a number: ie: N- or N,")
var ErrBadMinuteField error = errors.New("Bad minute field. The value received is not a number between 1-59")
var ErrBadHourField error = errors.New("Bad hour field. The value received is not a number between 1-23")
var ErrBadDomField error = errors.New("Bad day of month field. The value received is not a number between 1-31")
var ErrBadMonthField error = errors.New("Bad month field. The value received is not a number between 1-12")
var ErrBadDowField error = errors.New("Bad day of week field. The value received is not a number between 0-7")
var ErrBadRangeStep error = errors.New("Bad step provided for range syntax. This means that the step is either smaller or larger than the lower and upper bounds of the time part where it was found or the step number could not be parsed or the syntax was incorrect and more than one '/' was found")

func all(listOfCond []bool, value bool) bool {
	for _, c := range listOfCond {
		if c != value {
			return false
		}
	}
	return true
}

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

func addOrdinalSuffix(input string) string {
	switch input[len(input)-1] {
	case '1':
		return input + "st"
	case '2':
		return input + "nd"
	case '3':
		return input + "rd"
	default:
		return input + "th"
	}
}

func getPrettyMonthName(c *Cron, natMonthVal int64) string {
	return c.monthNames[natMonthVal-1]
}

func parseTemplate(tpl string, opts interface{}) string {
	var buf bytes.Buffer
	tplActual := template.Must(template.New("t").Parse(tpl))
	tplActual.Execute(&buf, opts)
	return buf.String()
}

type component struct {
	all        bool
	isRange    bool
	isStep     bool
	isList     bool
	isSingular bool
	skip       bool
	start      int64
	stop       int64
	values     []int64
	// override can be used when the start and stop values wouldn't make sense
	// for output like in the case of using a list of values
	override string
}

type parsedOutput struct {
	minutes []int64
	hours   []int64
	dom     []int64
	month   []int64
	dow     []int64
}

type rawParts struct {
	minutes string
	hours   string
	dom     string
	month   string
	dow     string
}

type Cron struct {
	monthNames [12]string
	dowNames   [8]string
	joiners    map[string]string
	rawParts
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
		[8]string{
			"sunday",
			"monday",
			"tueday",
			"wednesday",
			"thursday",
			"friday",
			"saturday",
			"sunday",
		},
		map[string]string{
			"hours": "of",
			"dom":   "on the",
			"month": "of",
			"dow":   "on",
		},
		rawParts{},
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
	c.rawParts.minutes = rawMinute
	c.rawParts.hours = rawHour
	c.rawParts.dom = rawDom
	c.rawParts.month = rawMonth
	c.rawParts.dow = rawDow

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
				if rawDow == v[0:3] {
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

	// Parse into components
	minComp := component{
		len(parsed.minutes) == 60,
		strings.Contains(c.rawParts.minutes, "-"),
		strings.Contains(c.rawParts.minutes, "/"),
		strings.Contains(c.rawParts.minutes, ","),
		false,
		false,
		parsed.minutes[0],
		parsed.minutes[len(parsed.minutes)-1],
		parsed.minutes,
		"",
	}
	minComp.isSingular = all(
		[]bool{
			minComp.isRange,
			minComp.isStep,
			minComp.isList,
			c.rawParts.minutes == "*",
		}, false)

	hourComp := component{
		len(parsed.hours) == 24,
		strings.Contains(c.rawParts.hours, "-"),
		strings.Contains(c.rawParts.hours, "/"),
		strings.Contains(c.rawParts.hours, ","),
		false,
		false,
		parsed.hours[0],
		parsed.hours[len(parsed.hours)-1],
		parsed.hours,
		"",
	}
	hourComp.isSingular = all(
		[]bool{
			hourComp.isRange,
			hourComp.isStep,
			hourComp.isList,
			c.rawParts.hours == "*",
		}, false)

	domComp := component{
		len(parsed.dom) == 31,
		strings.Contains(c.rawParts.dom, "-"),
		strings.Contains(c.rawParts.dom, "/"),
		strings.Contains(c.rawParts.dom, ","),
		false,
		false,
		parsed.dom[0],
		parsed.dom[len(parsed.dom)-1],
		parsed.dom,
		"",
	}
	domComp.isSingular = all(
		[]bool{
			domComp.isRange,
			domComp.isStep,
			domComp.isList,
			c.rawParts.dom == "*",
		}, false)

	monthComp := component{
		len(parsed.month) == 12,
		strings.Contains(c.rawParts.month, "-"),
		strings.Contains(c.rawParts.month, "/"),
		strings.Contains(c.rawParts.month, ","),
		false,
		false,
		parsed.month[0],
		parsed.month[len(parsed.month)-1],
		parsed.month,
		"",
	}
	monthComp.isSingular = all([]bool{
		monthComp.isRange,
		monthComp.isStep,
		monthComp.isList,
		c.rawParts.month == "*",
	}, false)

	dowComp := component{
		len(parsed.dow) == 8,
		strings.Contains(c.rawParts.dow, "-"),
		strings.Contains(c.rawParts.dow, "/"),
		strings.Contains(c.rawParts.dow, ","),
		false,
		false,
		parsed.dow[0],
		parsed.dow[len(parsed.dow)-1],
		parsed.dow,
		"",
	}
	dowComp.isSingular = all([]bool{
		dowComp.isRange,
		dowComp.isStep,
		dowComp.isList,
		c.rawParts.dow == "*",
	}, false)

	// The 5 fields of a cron can be broken down into 3 sets of components:
	//
	//  Time Component  , which is minutes and hours
	//  Day Component   , which is day-of-month and day-of-week
	//  Month Component , the month
	//
	// We'll process the parsed components in that order (Time, Day, Month) and
	// build the string from left to right in terms of said order.
	//
	// Specificity, when components are set to "all" we can ignore them in the
	// output and let the human reader infer their values

	// Quick Exit, if we have all stars then we'll just return the minutes portition
	if func() bool {
		for _, c := range []bool{minComp.all, hourComp.all, domComp.all, monthComp.all, dowComp.all} {
			if c == false {
				return false
			}
		}
		return true
	}() {
		return "every minute", nil
	}

	tpl := `{{.TimeComponent}}{{.DayComponent}}{{.MonthComponent}}`

	// Time Component
	// --------------------------------------------------------------------------
	tcTpl := func() string {
		var tcTpl string
		// minute
		if minComp.all {
			tcTpl += "every minute "
		}

		if minComp.isRange {
			tcTpl = `on minutes {{.MinStart}} through {{.MinStop}} `
		}

		if minComp.isList {
			ln := len(minComp.values) - 1
			last := minComp.values[ln]

			beforeLast := func(a []int64) string {
				var temp []string
				for _, i := range a {
					temp = append(temp, strconv.Itoa(int(i)))
				}
				return strings.Join(temp, ", ")
			}(minComp.values[0:ln])

			minComp.override = fmt.Sprintf("%s, and %d", beforeLast, last)
			tcTpl = `on minutes {{.MinOverride}} `
		}

		if minComp.isSingular {
			tcTpl = `at minute {{.MinStart}} `
		}

		if minComp.isStep {
			temp := strings.Split(c.rawParts.minutes, "/")
			minComp.override = temp[len(temp)-1]
			tcTpl = `every {{.MinOverride}} minutes `
		}

		// hour
		if hourComp.all {
			return strings.TrimRight(tcTpl, " ")
		}

		if hourComp.isSingular {
			tcTpl += `past {{.HourStart}}`
			return tcTpl
		}

		if hourComp.isStep {
			temp := strings.Split(c.rawParts.hours, "/")
			hourComp.override = temp[len(temp)-1]
			tcTpl += `past every {{.HourOverride}} hours`
			return tcTpl
		}

		if !hourComp.isSingular {
			tcTpl += `past the hours of `

			if hourComp.isList {
				ln := len(hourComp.values) - 1
				last := hourComp.values[ln]

				beforeLast := func(a []int64) string {
					var temp []string
					for _, i := range a {
						temp = append(temp, strconv.Itoa(int(i)))
					}
					return strings.Join(temp, ", ")
				}(hourComp.values[0:ln])

				hourComp.override = fmt.Sprintf("%s, and %d", beforeLast, last)
				tcTpl += `{{.HourOverride}}`
			} else {
				tcTpl += `{{.HourStart}} through {{.HourStop}}`
			}

		}
		return tcTpl
	}()
	tcRendered := parseTemplate(tcTpl, struct {
		MinStart     interface{}
		MinStop      interface{}
		MinOverride  interface{}
		HourStart    interface{}
		HourStop     interface{}
		HourOverride interface{}
	}{
		minComp.start,
		minComp.stop,
		minComp.override,
		hourComp.start,
		hourComp.stop,
		hourComp.override,
	})

	// Day Component
	// --------------------------------------------------------------------------
	dcTpl := func() string {
		if domComp.all && dowComp.all {
			return ""
		}

		var dcTpl string

		if domComp.isStep {
			dcTpl = " every "
		}

		if domComp.isRange {
			start := addOrdinalSuffix(strconv.Itoa(int(domComp.start)))
			stop := addOrdinalSuffix(strconv.Itoa(int(domComp.stop)))
			domComp.override = fmt.Sprintf(" on the %s through the %s", start, stop)
			dcTpl += "{{.DayOverride}}"
		}

		if domComp.isList {
			ln := len(domComp.values) - 1
			last := domComp.values[ln]

			beforeLast := func(a []int64) string {
				var temp []string
				for _, i := range a {
					temp = append(temp, addOrdinalSuffix(strconv.Itoa(int(i))))
				}
				return strings.Join(temp, ", ")
			}(domComp.values[0:ln])

			domComp.override = fmt.Sprintf("%s and the %s", beforeLast, addOrdinalSuffix(strconv.Itoa(int(last))))
			dcTpl += " on the {{.DayOverride}}"
		}

		if domComp.isSingular {
			domComp.override = addOrdinalSuffix(strconv.Itoa(int(domComp.start)))
			dcTpl += " on the {{.DayOverride}}"
		}

		if domComp.isStep {
			temp := strings.Split(c.rawParts.dom, "/")
			domComp.override = temp[len(temp)-1]
			dcTpl += "{{.DayOverride}} days"
		}

		// dow
		//
		// Given that the dow is tied to the dom in terms of joining phrases we
		// need to see if we have to check the dom here to see what if any we
		// should use
		var joiner string
		if domComp.all {
			joiner = " on"
		} else {
			joiner = " and on"
		}

		if dowComp.isRange {
			start := strings.Title(c.dowNames[dowComp.start])
			stop := strings.Title(c.dowNames[dowComp.stop])
			dowComp.override = fmt.Sprintf("%s %s through %s", joiner, start, stop)
			dcTpl += "{{.WeekDayOverride}}"
		}

		if dowComp.isList {
			ln := len(dowComp.values) - 1
			last := dowComp.values[ln]

			beforeLast := func(a []int64) string {
				var temp []string
				for _, i := range a {
					temp = append(temp, fmt.Sprintf("%ss", strings.Title(c.dowNames[i])))
				}
				return strings.Join(temp, ", ")
			}(dowComp.values[0:ln])

			dowComp.override = fmt.Sprintf("%s %s, and %ss", joiner, beforeLast, strings.Title(c.dowNames[last]))
			dcTpl += "{{.WeekDayOverride}}"
		}

		if dowComp.isSingular {
			dowComp.override = fmt.Sprintf("%s %ss", joiner, strings.Title(c.dowNames[dowComp.start]))
			dcTpl += "{{.WeekDayOverride}}"
		}

		if dowComp.isStep {
			temp := strings.Split(c.rawParts.dom, "/")
			dowComp.override = fmt.Sprintf("%s every %s days of the week", joiner, temp[len(temp)-1])
			dcTpl += "{{.WeekDayOverride}}"
		}

		return dcTpl
	}()
	dcRendered := parseTemplate(dcTpl, struct {
		DayStart        interface{}
		DayStop         interface{}
		DayOverride     interface{}
		WeekDayStart    interface{}
		WeekDayStop     interface{}
		WeekDayOverride interface{}
	}{
		domComp.start,
		domComp.stop,
		domComp.override,
		dowComp.start,
		dowComp.stop,
		dowComp.override,
	})

	// Month Component
	// --------------------------------------------------------------------------
	mcTpl := func() string {
		var mcTpl string

		// When it comes to the month we only care to add that to the output when
		// it's specific, so if its sets to all we can call it good
		if monthComp.all {
			return mcTpl
		}

		if monthComp.isSingular {
			mcTpl = ` of {{.MonthStart}}`
		}

		if monthComp.isRange {
			mcTpl = ` of {{.MonthStart}} through {{ .MonthStop }}`
		}

		// Both of these conditions produces a range of months so the template
		// should be the same
		if monthComp.isList || monthComp.isStep {
			prettyMonthNames := func() string {
				var res []string
				for i, v := range monthComp.values {
					if i+1 != len(monthComp.values) {
						res = append(res, strings.Title(getPrettyMonthName(c, v)+", "))
					} else {
						res = append(res, "and "+strings.Title(getPrettyMonthName(c, v)))
					}

				}
				return strings.Join(res, "")
			}()

			mcTpl = fmt.Sprintf(` of %s`, prettyMonthNames)
		}

		return mcTpl
	}()
	mcRendered := parseTemplate(mcTpl, struct {
		MonthStart    interface{}
		MonthStop     interface{}
		MonthOverride interface{}
	}{
		strings.Title(getPrettyMonthName(c, monthComp.start)),
		strings.Title(getPrettyMonthName(c, monthComp.stop)),
		monthComp.override,
	})

	// Finalize
	// --------------------------------------------------------------------------
	rtn := parseTemplate(tpl, struct {
		TimeComponent  string
		DayComponent   string
		MonthComponent string
	}{
		tcRendered,
		dcRendered,
		mcRendered,
	})

	return rtn, nil

}
