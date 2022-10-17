package ceres

import (
    "regexp"
    "fmt"
)

var SPLITTER = regexp.MustCompile("[ ,\t']")

func (c *CERES) ParseSentence(sentence string) []Entity {
    var tokens []Word

    pre_tokens := SPLITTER.Split(sentence, len(sentence))

    tokens = make([]Word, 0, len(pre_tokens))

    for _, token := range pre_tokens {
        if len(token) > 0 {
            tokens = append(tokens, Word(token).Lower())
        }
    }

    // the sentence is split in words
    fmt.Println("Now solving pronouns.")

    var tokens_used []bool= make([]bool, len(tokens), len(tokens))

    // checking for pronouns
    for i, word := range tokens {
        if c.pcs.IsPronoun(word) {
            tokens_used[i] = true
        }
    }

    return nil
}
