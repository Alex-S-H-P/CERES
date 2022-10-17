package ceres

import (
    "strings"
)

type Word string

func (w Word) Lower() Word{
    return Word(strings.ToLower(string(w)))
}
