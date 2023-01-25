package ceres

import (
	"fmt"
	"sync"
)

type UCS struct {
	unrecognized_words []Word
	uwm                sync.Mutex

	ceres_main *EntityType
}

func (ucs *UCS) Initialize() {
	ucs.unrecognized_words = make([]Word, 64)
}

/*
Creates an entity that is not based on any words.
*/
func (ucs *UCS) makeNonWordEntity(token string) *RecognizedEntity {
	// TODO handle all cases.
	switch recognizeType(token) {
	case TOKEN_TYPE_PRIC:
		// make a temporary price entity
	case TOKEN_TYPE_CURR:
		// replace with the currency's entity.
	case TOKEN_TYPE_INTN, TOKEN_TYPE_NUMB:
		// make temporary number entity
	case TOKEN_TYPE_UNKN:
		// fetch unknown. If failed make one.
	}

	// FIXME: this code is temporary
	ei := new(EntityInstance)

	return &RecognizedEntity{entity: ei,
		possessive: false, attribute: false, proposer: ucs}
}

/*
Due  to the specificicity of this proposer, it always returns a slice of only one element

Also, adds the element to the list of unrecognized_words
*/
func (ucs *UCS) proposeOptions(w Word, ctx *CTX) []*RecognizedEntity {
	fmt.Println("Adding unrecognized element", w, "to the UCS list")

	re := ucs.makeNonWordEntity(string(w))

	ucs.uwm.Lock()
	ucs.unrecognized_words = append(ucs.unrecognized_words, w)
	ucs.uwm.Unlock()

	return []*RecognizedEntity{re}
}

func (ucs *UCS) computeP(re *RecognizedEntity, ctx *CTX, previous ...*RecognizedEntity) float64 {
	// TODO: detect key-words related to definitions and setting values.
	return 0.1
}

func (ucs *UCS) name() string { return "UCS" }
