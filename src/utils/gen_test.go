package utils

import (
    "testing"
    "fmt"
)

func TestGen(t *testing.T) {
    mg := MapGenerator[int, string]{}
    fmt.Println(Generator[int](&mg))
}
