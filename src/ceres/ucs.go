package ceres

import (
    "fmt"
    "sync"
)

type UCS struct {
    unrecognized_words []Word
    uwm sync.Mutex

    ceres_main *EntityType

}

func (ucs*UCS) Initialize() {
    ucs.unrecognized_words = make([]Word, 64)
}

/*
Due  to the specificicity of this proposer, it always returns a slice of only one element

Also, adds the element to the list of unrecognized_words
*/
func (ucs*UCS)proposeOptions(w Word, ctx*CTX) []RecognizedEntity {
    fmt.Println("Adding unrecognized element", w, "to the UCS list")

    ei := new(EntityInstance)
    ei.typeOf = ucs.ceres_main
    re := MakeRecognizedEntity(ei, false, false, ucs, string(w))

    ucs.uwm.Lock()
    ucs.unrecognized_words = append(ucs.unrecognized_words, w)
    ucs.uwm.Unlock()

    return []RecognizedEntity{re}
}

func (ucs*UCS)computeP(re RecognizedEntity, ctx *CTX, previous...RecognizedEntity) float64{
    // TODO: detect key-words related to definitions and setting values.
    return 0.5
}
