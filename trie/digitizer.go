package trie

import "strings"

// Digitizer
type Digitizer interface {

    // Base returns the base for the Digitizer.
    Base() int

    // IsPrefixFree returns true if and only if the Digitizer guarantees that no element is a prefix of another.
    IsPrefixFree() bool

    // NumDigitsOf returns the number of digits in the provided element.
    NumDigitsOf(element interface{}) int

    // DigitOf returns the element of digit place for the provided element.
    DigitOf(element interface{}, place int) int

    // FormatDigit returns a string representation of the digit in the place specified for the given element.
    FormatDigit(element interface{}, place int) string
}

type stringDigitizer struct {
    base int
}

// NewStringDigitizer creates a new Digitizer for strings with the provided alphabet size.
func NewStringDigitizer(alphabetSize int) Digitizer {
    return &stringDigitizer{ base: alphabetSize + 1 }
}

// Base the base of the alphabet that includes the end of string character.
func (d *stringDigitizer) Base() int {
    return d.base
}

// IsPrefixFree returns true since this is a prefix free digitizer.
func (d *stringDigitizer) IsPrefixFree() bool {
    return true
}

// NumDigitsOf returns the number of digits in the provided string including the end of string character.
func (d *stringDigitizer) NumDigitsOf(element interface{}) int {
    return len(element.(string)) + 1
}

// DigitOf returns the integer element mapped to by the digit in the given place.
func (d *stringDigitizer) DigitOf(element interface{}, place int) int {
    if place >= len(element.(string)) {
        return 0
    } else {
        return int(strings.ToLower(element.(string))[place] - 'a' + 1)
    }
}

// FormatDigit returns a string representation of the digit in the place specified for the given element where '#' is
// used for the end of string character.
func (d *stringDigitizer) FormatDigit(element interface{}, place int) string {
    if place >= len(element.(string)) {
        return "#"
    } else {
        return string(strings.ToLower(element.(string))[place])
    }
}