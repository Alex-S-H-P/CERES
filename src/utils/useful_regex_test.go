package utils

import (
    "testing"
    re "regexp"
)

var patterns_to_test = [...]string{IntNumberPattern, NumberPattern, CurrencyPattern,
    PricePattern, WordPattern}
var testable_strings = [...]string{"€", "13", "hello", "1.25 €", "0.0", "'m", "#"}

var coherence = [len(patterns_to_test)][len(testable_strings)]bool{{false,
    true, false, false, false, false, false},
    {false, true, false, false, true, false, false},
    {true, false, false, false, false, false, false},
    {false, false, false, true, false, false, false},
    {false, false, true, false, false, true, false}}

func TestRegex(t *testing.T) {
    for i, pattern_to_test := range patterns_to_test {
        for j, test_string := range testable_strings {
            ok, err := re.MatchString("^"+pattern_to_test+"$", test_string)
            if err != nil {
                t.Errorf("Pattern %v (n°%v) is not valid", pattern_to_test, i)
            }
            if (ok != coherence[i][j]) {
                shouldq := "should"
                if ok {
                    shouldq += "n't"
                }
                t.Errorf("Pattern %v (n°%v) and sentence %v, (n°%v) %s match...", pattern_to_test, i, test_string, j, shouldq)
            }
        }
    }
}
