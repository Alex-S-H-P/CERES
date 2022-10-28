package ceres

import (
    "testing"
    "fmt"
)

func P_Test(t *testing.T) {
    var p = new(Pronoun)
    genders := []int8 {MALE, FEMALE, NEUTRAL, UNKNOWN}
    number  := []int8 {SINGULAR, PLURAL, DUAL}
    person  := []int8 {PERSON1, PERSON2, PERSON3}
    for _, gender := range genders {
        for _, number := range number {
            for _, person := range person {
                p.MakeGNP(gender, number, person)
                g, n, p := p.GNP_Sep()
                if g == gender && n == number && p == person {
                    t.Error(fmt.Errorf("Mismatch between (%v, %v, %v) != (%v, %v, %v)",
                        g, n, p, gender, number, person))
                }
            }
        }
    }
}
