## Number

Breaks up a contiguous number.

```gherkin
Feature: number

	Scenario: Take a number and output in words
		Given <input> like "100000000"
		When ran like
			`human number -w <input>`
		Then it should output
			"one hundred million"

	Scenario: Take a number and output visual separators
		Given <input> like "100000000"
		When ran like
			`human number -g <input>`
		Then it should output
			"100,000,000"
```

## Size

Translates byte sizes into size acronym

```gherkin
Feature: size

	Scenario: Take a number that represents bytes and give back nearest acronym
		Given <input> like "5000000000"
		When ran like
			`human size <input>`
		Then it should output
			"5G"

	Scenario:  Take a size in computer acronym, and output bytes
		Given <input> like "5G" or "5 gigs"
		When ran like
			`human size <input>`
		Then it should output
			"5000000000"
```
