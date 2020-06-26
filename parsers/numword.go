// Given a number in string form
// Return the number in word form of the greatest power (e.g. 300,100,000,000 => 300.1 Billion)
// Lower limit > Thousands (e.g. xxx,000) # Numbers below this are in hundreds and can be expressed numerically
// Upper limit > Centillion (10^303) ref: https://en.wikipedia.org/wiki/Names_of_large_numbers#Standard_dictionary_numbers

package parsers

import (
  "regexp"
)

// NumberWord handles strings made of contiguous "0-9" characters
// converts to word strings of the greatest power
type NumberWord struct{}

// NewNumberWord construcs a NumberWord struct
func NewNumberWord() *NumberWord {
  return &NumberWord{}
}

// CanParseIntoHuman ...
func (n *NumberWord) CanParseIntoHuman(s string) bool {
  // is it 4 or more characters? (e.g. is it => 1000)
  if len(s) < 4 {
    return false
  }

  // is it a (delimited[,. ]) number?
  match, _ := regexp.MatchString(`^(([0-9]+)|([0-9]{1,3}[., ])+[0-9]{1,3})$`, s)
  if match {
    // is it less than the max? 
      //  300000 max places plus delimiters
    if len(s) < 400000 {
      return true
    }
  }
  // everything else is not a number
  return false
}

// CanParseFromHuman ...
func (n *NumberWord) CanParseFromHuman(s string) bool {
  // is it a digit word combo? (e.g. 48 billion) /\d+ [a-zA-Z]+)
  // is the word in the translation map?
  return false
}

// DoIntoHuman ...
func (n *NumberWord) DoIntoHuman(s string) string {
  // Create a NumberGroup from string
  // Split the NumberGroup string into an array ng[]
  // Compare the len of ng[] with numwords translation index
  // Round ng[1] to hundreds place (e.g. 155 => 200) = decimals
  // Return ng[0]+ "." + decimal  + " " + numwords[len(ng)] 
  return "100.0 FOOillion"
}

// DoFromHuman ...
// Only works with highest power (e.g. 100.3 Billion, not 100,300 Million)
func (n *NumberWord) DoFromHuman(s string) string {
  // Split numbers from word
  // Get powers from translation map
  // Return ( numbers * 10^foo )
  return "100"
}
